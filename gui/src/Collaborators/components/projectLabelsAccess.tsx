import * as React from 'react';
import * as _ from 'lodash';
import Table from '../../Projects/components/Table';
import Row from '../../Projects/components/Row';
import Cell from '../../Projects/components/Cell';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import '../styles/collaborators.scss';
import { fetchLabelsForProject } from '../actions/collaborators.actions';
import { setCurrentProject } from '../../Projects/actions/projects.actions';

interface Props {
  projectid: string,
  labels: Array<any>
}

interface DispatchProps {
  fetchLabelsForProject: Function,
  setCurrentProject: Function
}

export class ProjectLabelsAccess extends React.Component<Props & DispatchProps, any> {

  componentWillMount(): void {
    this.props.setCurrentProject(parseInt(this.props.projectid, 10));
    this.props.fetchLabelsForProject();
  }

  render(): React.ReactElement<HTMLDivElement> {
    return (
      <div className="labelsAccess">
        <p></p>
        <h1>Labels Access</h1>
        <p>A label is applied to a particular model to designate it for particular use, e.g. a 'prod' label to make a model as fit for production. Labels have restricted access control to make sure only users with appropriate privileges can change which models have a particular label.</p>
        <p>All labels for this project, along with users privileges for that label, are listed below.</p>
        <Table>
          <Row header={true}>
            <Cell>LABEL</Cell>
            <Cell>USERS</Cell>
          </Row>
            {this.props.labels ?
              this.props.labels.map((label, labelIndex) => {
                return <Row key={labelIndex}>
                  <Cell>{label.name}</Cell>
                  <Cell>
                  {label.identities ? label.identities.map((identity, identityIndex) => {
                    return <div key={identityIndex}>
                      <span className="access-name">{identity.identity_name}</span>&nbsp;
                      <span className="access-type">{identity.kind}</span>
                    </div>;
                  }) : null}
                  </Cell>
                </Row>;
              }) : null}
        </Table>
      </div>
    );
  }
}

function mapStateToProps(state) {
  return {
    labels: state.collaborators.labels
  };
}

function mapDispatchToProps(dispatch) {
  return {
    fetchLabelsForProject: bindActionCreators(fetchLabelsForProject, dispatch),
    setCurrentProject: bindActionCreators(setCurrentProject, dispatch)
  };
}

export default connect<any, DispatchProps, any>(mapStateToProps, mapDispatchToProps)(ProjectLabelsAccess);
