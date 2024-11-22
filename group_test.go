package main

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var spaceRe = regexp.MustCompile(`\s+`)
var lowerSnakeRe = regexp.MustCompile(`^[a-z]+(?:_[_a-z]+)?$`)

func checkLowerSnake(t *testing.T, name, value string) {
	t.Helper()

	if !lowerSnakeRe.MatchString(value) {
		t.Errorf("%s is %q, want match for %q", name, value, lowerSnakeRe.String())
	}
}

func TestGroupFieldNames(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input *group
		want  []string
	}{
		{
			name:  "empty",
			input: &group{},
		},
		{
			name: "one each",
			input: &group{
				keyFields: []*descriptor{
					{fieldName: "keyField"},
				},
				infoFields: []*descriptor{
					{fieldName: "infoField"},
				},
				metricFields: []*descriptor{
					{fieldName: "metricField"},
				},
			},
			want: []string{"keyField", "infoField", "metricField"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.input.fieldNames()

			if diff := cmp.Diff(got, tc.want, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("fieldNames() difference (-got +want):\n%s", diff)
			}
		})
	}
}

func checkSingleDescriptor(t *testing.T, g *group, d *descriptor) {
	t.Helper()

	checkLowerSnake(t, "fieldName", d.fieldName)
	checkLowerSnake(t, "metricName", d.metricName)

	if wantPrefix := fmt.Sprintf("%s_", g.name); !strings.HasPrefix(d.metricName, wantPrefix) {
		t.Errorf("metricName %q must start with %q", d.metricName, wantPrefix)
	} else if strings.HasPrefix(d.metricName, wantPrefix+wantPrefix) {
		t.Errorf("metricName %q has double %q prefix", d.metricName, wantPrefix)
	}

	wantDesc := spaceRe.ReplaceAllLiteralString(d.desc, " ")
	wantDesc = strings.TrimSpace(wantDesc)
	wantDesc = strings.TrimSuffix(wantDesc, ".")

	if diff := cmp.Diff(d.desc, wantDesc); diff != "" {
		t.Errorf("Description not normalized (-got +want):\n%s", diff)
	}
}

func checkReportFields(t *testing.T, g *group, fields []*descriptor) {
	t.Helper()

	sortedFields := append([]*descriptor(nil), fields...)

	slices.SortFunc(sortedFields, func(a, b *descriptor) int {
		return strings.Compare(a.fieldName, b.fieldName)
	})

	if diff := cmp.Diff(fields, sortedFields, cmp.AllowUnexported(descriptor{}), cmpopts.IgnoreTypes(metricValueFunc(nil))); diff != "" {
		t.Errorf("descriptors not sorted (-got +want):\n%s", diff)
	}

	for idx, d := range fields {
		t.Run(fmt.Sprintf("%d:%s", idx, d.fieldName), func(t *testing.T) {
			checkSingleDescriptor(t, g, d)
		})
	}
}

func checkReportDescriptors(t *testing.T, g *group) {
	t.Helper()

	fieldNames := []string{}
	fieldNameMap := map[string]struct{}{}
	metricNameMap := map[string]struct{}{}

	for _, tc := range []struct {
		name   string
		fields []*descriptor
	}{
		{"key", g.keyFields},
		{"info", g.infoFields},
		{"metric", g.metricFields},
	} {
		t.Run(tc.name, func(t *testing.T) {
			checkReportFields(t, g, tc.fields)
		})

		for _, d := range tc.fields {
			if _, ok := fieldNameMap[d.fieldName]; ok {
				t.Errorf("Duplicate field name %q", d.fieldName)
			}

			if _, ok := metricNameMap[d.metricName]; ok {
				t.Errorf("Duplicate metric name %q", d.metricName)
			}

			fieldNames = append(fieldNames, d.fieldName)
			fieldNameMap[d.fieldName] = struct{}{}
			metricNameMap[d.metricName] = struct{}{}
		}
	}

	if diff := cmp.Diff(g.fieldNames(), fieldNames, cmpopts.SortSlices(func(a, b string) bool {
		return a < b
	})); diff != "" {
		t.Errorf("fieldNames() difference (-got +want):\n%s", diff)
	}
}

func TestGroupValidation(t *testing.T) {
	kindRe := regexp.MustCompile(`^[a-z]{2,6}$`)

	for _, g := range allGroups {
		t.Run(g.name.String(), func(t *testing.T) {
			if !kindRe.MatchString(g.name.String()) {
				t.Errorf("report kind is %q, want match for %q", g.name.String(), kindRe.String())
			}

			checkLowerSnake(t, "infoMetricName", g.infoMetricName)
			checkReportDescriptors(t, g)
		})
	}
}
