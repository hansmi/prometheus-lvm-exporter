package main

import (
	"fmt"
	"regexp"
	"slices"
	"sort"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var spaceRe = regexp.MustCompile(`\s+`)
var lowerSnakeRe = regexp.MustCompile(`^[a-z]+(?:_[_a-z]+)?$`)

func checkLowerSnake(t *testing.T, name, value string) {
	t.Helper()

	if !lowerSnakeRe.MatchString(value) {
		t.Errorf("%s is %q, want match for %q", name, value, lowerSnakeRe.String())
	}
}

func checkSingleField(t *testing.T, g *group, f field) {
	t.Helper()

	checkLowerSnake(t, "Name", f.Name())
	checkLowerSnake(t, "MetricName", f.MetricName())

	if wantPrefix := fmt.Sprintf("%s_", g.name); !strings.HasPrefix(f.MetricName(), wantPrefix) {
		t.Errorf("Metric name %q must start with %q", f.MetricName(), wantPrefix)
	} else if strings.HasPrefix(f.MetricName(), wantPrefix+wantPrefix) {
		t.Errorf("Metric name %q has double %q prefix", f.MetricName(), wantPrefix)
	}

	wantHelp := spaceRe.ReplaceAllLiteralString(f.Help(), " ")
	wantHelp = strings.TrimSpace(wantHelp)
	wantHelp = strings.TrimSuffix(wantHelp, ".")

	if diff := cmp.Diff(f.Help(), wantHelp); diff != "" {
		t.Errorf("Help text not normalized (-got +want):\n%s", diff)
	}
}

func checkGroupFields[T field](t *testing.T, g *group, add func(*testing.T, field), fields []T) {
	t.Helper()

	var names []string

	for idx, f := range fields {
		names = append(names, f.Name())

		t.Run(fmt.Sprintf("%d:%s", idx, f.Name()), func(t *testing.T) {
			add(t, f)
			checkSingleField(t, g, f)
		})
	}

	sortedNames := slices.Clone(names)

	sort.Strings(sortedNames)

	if diff := cmp.Diff(names, sortedNames); diff != "" {
		t.Errorf("Fields not sorted (-got +want):\n%s", diff)
	}
}

func checkGroup(t *testing.T, g *group) {
	t.Helper()

	fieldNameMap := map[string]struct{}{}
	metricNameMap := map[string]struct{}{}

	checkLowerSnake(t, "infoMetricName", g.infoMetricName)
	metricNameMap[g.infoMetricName] = struct{}{}

	process := func(t *testing.T, f field) {
		if _, ok := fieldNameMap[f.Name()]; ok {
			t.Errorf("Duplicate field name %q", f.Name())
		}

		if _, ok := metricNameMap[f.MetricName()]; ok {
			t.Errorf("Duplicate metric name %q", f.MetricName())
		}

		fieldNameMap[f.Name()] = struct{}{}
		metricNameMap[f.MetricName()] = struct{}{}
	}

	checkGroupFields(t, g, process, g.keyFields)
	checkGroupFields(t, g, process, g.textFields)
	checkGroupFields(t, g, process, g.numericFields)
}

func TestGroupValidation(t *testing.T) {
	kindRe := regexp.MustCompile(`^[a-z]{2,6}$`)

	for _, g := range allGroups {
		t.Run(g.name.String(), func(t *testing.T) {
			if !kindRe.MatchString(g.name.String()) {
				t.Errorf("Group kind is %q, want match for %q", g.name.String(), kindRe.String())
			}

			checkGroup(t, g)
		})
	}
}
