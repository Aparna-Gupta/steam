/**
 * Created by justin on 6/17/16.
 */

import * as React from 'react';
import * as ReactDOM from 'react-dom';
import { Router, Route, IndexRoute, hashHistory } from 'react-router';
import { syncHistoryWithStore } from 'react-router-redux';
import { Store, createStore, applyMiddleware, compose } from 'redux';
import { Provider } from 'react-redux';
import thunk from 'redux-thunk';
import App from './App/App';
import Clusters from './Clusters/Clusters';
import Projects from './Projects/Projects';
import ProjectDetails from './ProjectDetails/ProjectDetails';
import { rootReducer } from './App/reducers/rootReducer';

import './variables.scss';
import 'bootstrap/dist/css/bootstrap.css';
import 'font-awesome/css/font-awesome.css';

const initialState = {};

const store: Store = createStore(
  rootReducer,
  initialState,
  compose(applyMiddleware(thunk), window.devToolsExtension ? window.devToolsExtension() : f => f)
);

let history: ReactRouterRedux.ReactRouterReduxHistory = syncHistoryWithStore(hashHistory, store);

ReactDOM.render(
  <Provider store={store}>
    <Router history={history}>
      <Route path="/" component={App}>
        <IndexRoute component={Clusters}/>
        <Route path="clusters" component={Clusters}>
          <IndexRoute component={Clusters}/>
        </Route>
        <Route path="projects" component={Projects}>
        </Route>
        <Route path="/projects/:id" component={ProjectDetails}>
        </Route>
      </Route>
    </Router>
  </Provider>,
  document.getElementById('app'));