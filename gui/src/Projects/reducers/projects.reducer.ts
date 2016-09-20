/**
 * Created by justin on 7/18/16.
 */
import * as _ from 'lodash';
import {
  RECEIVE_CLUSTERS, RECEIVE_MODELS, CREATE_PROJECT_COMPLETED, SET_CURRENT_PROJECT,
  RECEIVE_PROJECTS, RECEIVE_DATASETS_FROM_CLUSTER, RECEIVE_MODELS_FROM_PROJECT, RECEIVE_PROJECT, REQUEST_CLUSTERS,
  REQUEST_MODELS, REGISTER_CLUSTER_ERROR
} from '../actions/projects.actions';

let initialState = {
  clusters: [],
  models: [],
  project: {},
  availableProjects: null,
  isClusterFetchInProcess: false
};

export const projectsReducer = (state = initialState, action: any) => {
  switch (action.type) {
    case SET_CURRENT_PROJECT :
      if (state.project.hasOwnProperty("id")) {
        if ((state.project as any).id === action.projectId) {
          return state;
        }
      }
      let toReturn: any =  _.assign({}, state);
      toReturn.project = { id: action.projectId };
      return toReturn;
    case REQUEST_CLUSTERS:
      return _.assign({}, state, {
        isClusterFetchInProcess: true
      });
    case RECEIVE_CLUSTERS:
      return _.assign({}, state, {
        clusters: action.clusters,
        isClusterFetchInProcess: false,
        registerClusterError: null
      });
    case REQUEST_MODELS:
      return _.assign({}, state, {
        isModelFetchInProcess: true
      });
    case RECEIVE_MODELS:
      return _.assign({}, state, {
        models: action.models,
        isModelFetchInProcess: false
      });
    case RECEIVE_PROJECT:
      return _.assign({}, state, {
        project: action.project
      });
    case RECEIVE_MODELS_FROM_PROJECT:
      return _.assign({}, state, {
        models: action.models
      });
    case CREATE_PROJECT_COMPLETED:
      return _.assign({}, state, {
        project: action.project
      });
    case RECEIVE_PROJECTS:
      return _.assign({}, state, {
        availableProjects: action.projects
      });
    case RECEIVE_DATASETS_FROM_CLUSTER:
      return _.assign({}, state, {
        datasets: action.datasets
      });
    case REGISTER_CLUSTER_ERROR:
      return _.assign({}, state, {
        registerClusterError: action.message
      });
    default:
      return state;
  }
};
