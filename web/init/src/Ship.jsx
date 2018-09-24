import React from "react";
import { Provider } from "react-redux";
import RouteDecider from "./containers/RouteDecider";
import AppWrapper from "./containers/AppWrapper";
import { configureStore } from "./redux";
import PropTypes from "prop-types";

import "./scss/index.scss";
const bodyClass = "ship-init";

export class Ship extends React.Component {
  constructor() {
    super();
    this.state = {
      store: null
    }
  }
  static propTypes = {
    /** API endpoint for the Ship binary */
    apiEndpoint: PropTypes.string.isRequired,
    /**
     * Base path name for the internal Ship Init component router
     * */
    basePath: PropTypes.string,
    /**
     * Determines whether or not the Ship Init app will instantiate its own BrowserRouter
     * */
    routerEnabled: PropTypes.bool,
    /**
     * Determines whether or not the header will be shown
     * */
    headerEnabled: PropTypes.bool,
    /**
     * Parent history needed to handle internal routing, in the case of routerEnabled = false
     * */
    history: PropTypes.object
  }

  componentDidMount() {
    const { apiEndpoint } = this.props;
    // This fixes a bug regarding explicit hot reloading of reducers introduced in Redux v2.0.0
    const store = configureStore(apiEndpoint);
    this.setState({ store });
  }

  render() {
    const { history = null, headerEnabled = true, basePath = "", routerEnabled = true } = this.props;
    const { store } = this.state;

    if(!store) return <div></div>;

    return (
      <div id="ship-init-component">
        <Provider store={store}>
          <AppWrapper>
            <RouteDecider 
              headerEnabled={headerEnabled}
              routerEnabled={routerEnabled} 
              basePath={basePath}
              history={history}
            />
          </AppWrapper>
        </Provider>
      </div>
    )
  }
}

