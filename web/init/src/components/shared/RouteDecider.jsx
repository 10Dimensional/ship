import React from "react";
import PropTypes from "prop-types";
import { BrowserRouter, Router, Route, Switch } from "react-router-dom";
import isEmpty from "lodash/isEmpty";
import NavBar from "../../containers/Navbar";

import Loader from "./Loader";
import StepNumbers from "./StepNumbers";
import DetermineComponentForRoute from "../../containers/DetermineComponentForRoute";
import StepDone from "./StepDone";

const ShipRoutesWrapper = ({ routes, headerEnabled, basePath }) => (
  <div className="flex-column flex1">
    <div className="flex-column flex1 u-overflow--hidden u-position--relative">
      {!routes ?
        <div className="flex1 flex-column justifyContent--center alignItems--center">
          <Loader size="60" />
        </div>
        :
        <div className="u-minHeight--full u-minWidth--full flex-column flex1">
          {headerEnabled && <NavBar hideLinks={true} routes={routes} basePath={basePath} />}
          <StepNumbers basePath={basePath} steps={routes} />
          <div className="flex-1-auto flex-column u-overflow--auto">
            <Switch>
              {routes && routes.map((route) => (
                <Route
                  exact
                  key={route.id}
                  path={`${basePath}/${route.id}`}
                  render={() => <DetermineComponentForRoute
                    basePath={basePath}
                    routes={routes} 
                    routeId={route.id} 
                  />}
                />
              ))}
              <Route exact path="/" component={() => <div className="flex1 flex-column justifyContent--center alignItems--center"><Loader size="60" /></div> } />
              <Route exact path="/done" component={() =>  <StepDone />} />
            </Switch>
          </div>
        </div>
      }
    </div>
  </div>
)

export default class RouteDecider extends React.Component {
  static propTypes = {
    routerEnabled: PropTypes.bool,
    isDone: PropTypes.bool.isRequired,
    routes: PropTypes.arrayOf(
      PropTypes.shape({
         id: PropTypes.string,
         description: PropTypes.string,
         phase: PropTypes.string,
      })
    )
  }

  componentDidUpdate(lastProps) {
    const {
      routes,
      getHelmChartMetadata,
    } = this.props
    if (routes !== lastProps.routes && routes.length) {
      for (let i = 0; i < routes.length; i++) {
        if (routes[i].phase.includes("helm")) {
          getHelmChartMetadata();
          break;
        }
      }
    }
  }

  componentDidMount() {
    if (isEmpty(this.props.routes)) {
      this.props.getRoutes();
    }
  }

  render() {
    const {
      routes,
      isDone,
      routerEnabled,
      basePath,
      history,
      headerEnabled
    } = this.props;
    const routeProps = {
      routes,
      basePath,
      headerEnabled
    }
    return (
      <div className="u-minHeight--full u-minWidth--full flex-column flex1">
        { !routerEnabled ? 
          <Router history={history}>
            <ShipRoutesWrapper 
              {...routeProps}
            />
          </Router> :
          <BrowserRouter basename="/">
            <ShipRoutesWrapper 
              {...routeProps}
            />
          </BrowserRouter>
        }
      </div>
    );
  }
}
