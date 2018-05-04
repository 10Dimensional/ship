package render

import (
	"context"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/replicatedcom/ship/pkg/lifecycle/render/state"
	_ "github.com/replicatedcom/ship/pkg/lifecycle/render/test-cases"

	"github.com/go-kit/kit/log"
	"github.com/golang/mock/gomock"
	"github.com/replicatedcom/ship/pkg/api"
	"github.com/replicatedcom/ship/pkg/lifecycle/render/planner"
	"github.com/replicatedcom/ship/pkg/test-mocks/config"
	mockplanner "github.com/replicatedcom/ship/pkg/test-mocks/planner"
	"github.com/replicatedcom/ship/pkg/test-mocks/ui"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

type testcase struct {
	Name        string
	Metadata    api.ReleaseMetadata
	Spec        api.Spec
	ViperConfig map[string]interface{} `yaml:"viper_config"`
	Responses   map[string]string
	Expect      map[string]string
}

func TestRender(t *testing.T) {
	ctx := context.Background()

	renderer := &Renderer{
		Logger: log.NewNopLogger(),
	}

	tests := loadTestCases(t, filepath.Join("test-cases", "render-inline.yaml"))

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			mc := gomock.NewController(t)

			mockUI := ui.NewMockUi(mc)
			p := mockplanner.NewMockPlanner(mc)
			configResolver := config.NewMockResolver(mc)
			mockFS := afero.Afero{Fs: afero.NewMemMapFs()}

			release := &api.Release{Spec: test.Spec}
			renderer.Fs = mockFS
			renderer.UI = mockUI
			renderer.ConfigResolver = configResolver
			renderer.Planner = p
			renderer.StateManager = &state.StateManager{
				Logger: renderer.Logger,
			}

			func() {
				defer mc.Finish()

				configResolver.EXPECT().
					ResolveConfig(ctx, release, gomock.Any()).
					Return(test.ViperConfig, nil)

				p.EXPECT().
					Build(test.Spec.Assets.V1, test.Spec.Config.V1, gomock.Any(), test.ViperConfig).
					Return(planner.Plan{})

				p.EXPECT().
					Execute(ctx, planner.Plan{}).
					Return(nil)

				err := renderer.Execute(ctx, release, &api.Render{})
				assert.NoError(t, err)
			}()
		})
	}
}

func loadTestCases(t *testing.T, path string) []testcase {
	tests := make([]testcase, 1)
	contents, err := ioutil.ReadFile(path)
	assert.NoError(t, err)
	err = yaml.UnmarshalStrict(contents, &tests)
	assert.NoError(t, err)
	return tests
}
