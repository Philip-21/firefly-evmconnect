package metrics

import (
	"context"

	"github.com/gorilla/mux"
	"github.com/hyperledger/firefly-common/pkg/config"
	"github.com/hyperledger/firefly-common/pkg/httpserver"
	"github.com/hyperledger/firefly-evmconnect/internal/evmconfig"
	"github.com/hyperledger/firefly-transaction-manager/pkg/fftm"
)

type Manager struct {
	ctx               context.Context
	metricsManager    Metrics
	metricsServer     httpserver.HTTPServer
	metricsEnabled    bool
	metricsHandled    bool
	metricsServerDone chan error
	fftm.Manager
}

func (m *Manager) CreateMetricsMuxRouter() *mux.Router {
	r := mux.NewRouter()
	r.Path(config.GetString(evmconfig.MetricsPath)).Handler(m.metricsManager.HTTPHandler())
	return r
}

func (m *Manager) MetricsServer() {
	if m.metricsEnabled {
		go m.Start()
	} else if m.metricsHandled {
		go m.Close()
	}
}

func (m *Manager) initServices(ctx context.Context) (err error) {
	if m.metricsEnabled {
		m.metricsServer, err = httpserver.NewHTTPServer(ctx, "metrics", m.CreateMetricsMuxRouter(), m.metricsServerDone, evmconfig.MetricsConfig, evmconfig.CorsConfig)
		if err != nil {
			return err
		}
	}
	return nil
}
