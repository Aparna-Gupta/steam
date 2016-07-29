/**
 * Created by justin on 6/27/16.
 */
import * as React from 'react';
import * as _ from 'lodash';
import Panel from '../Projects/components/Panel';
import PageHeader from '../Projects/components/PageHeader';
import { fetchModelsFromCluster, fetchClusters, registerCluster } from '../Projects/actions/projects.actions';
import { bindActionCreators } from 'redux';
import { Cluster } from '../Proxy/Proxy';
import { connect } from 'react-redux';
import './styles/clusters.scss';

interface DispatchProps {
  fetchClusters: Function
  fetchModelsFromCluster: Function
  registerCluster: Function
}

interface Props {
  clusters: Cluster[]
}

export class Clusters extends React.Component<Props & DispatchProps, any> {
  componentWillMount(): void {
    if (_.isEmpty(this.props.clusters)) {
      this.props.fetchClusters();
    }
  }

  render(): React.ReactElement<HTMLDivElement> {
    if (!this.props.clusters) {
      return <div></div>;
    }
    return (
      <div className="clusters">
        <PageHeader>CLUSTERS</PageHeader>
        <div className="panel-container">
          {this.props.clusters.map((cluster, i) => {
            console.log(cluster);
            return (
              <Panel key={i}>
                <header>
                  <span><i className="fa fa-cubes"/> <a href={cluster.address} target="_blank" rel="noopener">{cluster.name} @ {cluster.address}</a></span>
                </header>
                <article>
                  <h3>
                    STATUS
                  </h3>
                  <h2 className="cluster-status">
                    {cluster.state === 'started' ? 'OK' : cluster.state}
                  </h2>
                </article>
              </Panel>
            );
          })}
        </div>
      </div>
    );
  }
}

function mapStateToProps(state): any {
  return {
    clusters: state.projects.clusters
  };
}

function mapDispatchToProps(dispatch): DispatchProps {
  return {
    fetchClusters: bindActionCreators(fetchClusters, dispatch),
    fetchModelsFromCluster: bindActionCreators(fetchModelsFromCluster, dispatch),
    registerCluster: bindActionCreators(registerCluster, dispatch)
  };
}

export default connect<any, DispatchProps, any>(mapStateToProps, mapDispatchToProps)(Clusters);
