package models_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/prymitive/karma/internal/config"
	"github.com/prymitive/karma/internal/models"
)

type annotationMapsTestCase struct {
	defaultHidden bool
	annotationMap map[string]string
	annotations   models.Annotations
	visible       []string
	hidden        []string
}

var annotationMapsTestCases = []annotationMapsTestCase{
	{
		annotationMap: map[string]string{
			"foo": "bar",
		},
		annotations: models.Annotations{
			models.Annotation{
				Name:    "foo",
				Value:   "bar",
				Visible: true,
				IsLink:  false,
			},
		},
	},
	{
		annotationMap: map[string]string{
			"foo": "http://localhost",
		},
		annotations: models.Annotations{
			models.Annotation{
				Name:    "foo",
				Value:   "http://localhost",
				Visible: true,
				IsLink:  true,
			},
		},
	},
	{
		annotationMap: map[string]string{
			"foo": "ftp://localhost",
		},
		annotations: models.Annotations{
			models.Annotation{
				Name:    "foo",
				Value:   "ftp://localhost",
				Visible: true,
				IsLink:  true,
			},
		},
	},
	{
		annotationMap: map[string]string{
			"foo": "https://localhost/xxx",
			"abc": "xyz",
		},
		annotations: models.Annotations{
			models.Annotation{
				Name:    "abc",
				Value:   "xyz",
				Visible: true,
				IsLink:  false,
			},
			models.Annotation{
				Name:    "foo",
				Value:   "https://localhost/xxx",
				Visible: true,
				IsLink:  true,
			},
		},
	},
	{
		annotationMap: map[string]string{
			"notLink": "https://some-links.domain.com/healthcheck in dev (job: blackbox) is not successfully probing via the blackbox prober. this could be due to the endpoint being offline, returning an invalid status code, taking too long to respond, etc.",
		},
		annotations: models.Annotations{
			models.Annotation{
				Name:    "notLink",
				Value:   "https://some-links.domain.com/healthcheck in dev (job: blackbox) is not successfully probing via the blackbox prober. this could be due to the endpoint being offline, returning an invalid status code, taking too long to respond, etc.",
				Visible: true,
				IsLink:  false,
			},
		},
	},
	{
		annotationMap: map[string]string{
			"notLink": "mailto:me@example.com",
		},
		annotations: models.Annotations{
			models.Annotation{
				Name:    "notLink",
				Value:   "mailto:me@example.com",
				Visible: true,
				IsLink:  false,
			},
		},
	},
	{
		defaultHidden: true,
		visible:       []string{"visible"},
		annotationMap: map[string]string{
			"hidden": "value",
		},
		annotations: models.Annotations{
			models.Annotation{
				Name:    "hidden",
				Value:   "value",
				Visible: false,
				IsLink:  false,
			},
		},
	},
	{
		defaultHidden: true,
		visible:       []string{"visible"},
		hidden:        []string{"hidden"},
		annotationMap: map[string]string{
			"visible": "value",
			"hidden":  "value",
			"default": "value",
		},
		annotations: models.Annotations{
			models.Annotation{
				Name:    "default",
				Value:   "value",
				Visible: false,
				IsLink:  false,
			},
			models.Annotation{
				Name:    "hidden",
				Value:   "value",
				Visible: false,
				IsLink:  false,
			},
			models.Annotation{
				Name:    "visible",
				Value:   "value",
				Visible: true,
				IsLink:  false,
			},
		},
	},
}

func TestAnnotationsFromMap(t *testing.T) {
	for _, testCase := range annotationMapsTestCases {
		config.Config.Annotations.Default.Hidden = testCase.defaultHidden
		config.Config.Annotations.Hidden = testCase.hidden
		config.Config.Annotations.Visible = testCase.visible
		result := models.AnnotationsFromMap(testCase.annotationMap)
		if !reflect.DeepEqual(testCase.annotations, result) {
			t.Errorf("AnnotationsFromMap result mismatch for map %v, expected %v got %v",
				testCase.annotationMap, testCase.annotations, result)
		}
	}
}

func TestAnnotationsSort(t *testing.T) {
	annotations := models.Annotations{
		models.Annotation{
			Name:    "bar",
			Value:   "abc",
			Visible: true,
			IsLink:  false,
		},
		models.Annotation{
			Name:    "xyz",
			Value:   "xyz",
			Visible: true,
			IsLink:  true,
		},
		models.Annotation{
			Name:    "abc",
			Value:   "bar",
			Visible: true,
			IsLink:  true,
		},
	}
	sort.Stable(annotations)
	if annotations[0].Name != "abc" {
		t.Errorf("Expected 'abc' to be first, got '%s'", annotations[0].Name)
	}
	if annotations[2].Name != "xyz" {
		t.Errorf("Expected 'xyz' to be last, got '%s'", annotations[2].Name)
	}
}
