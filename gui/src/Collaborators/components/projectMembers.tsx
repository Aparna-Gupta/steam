import * as React from 'react';
import Table from '../../Projects/components/Table';
import Row from '../../Projects/components/Row';
import Cell from '../../Projects/components/Cell';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import '../styles/collaborators.scss';
import { fetchMembers } from '../actions/collaborators.actions';

interface Props {
  members: Array<any>
}

interface DispatchProps {
  fetchMembers: Function
}

export class ProjectMembers extends React.Component<Props & DispatchProps, any> {
  componentWillMount(): void {
    this.props.fetchMembers();
  }

  render(): React.ReactElement<HTMLDivElement> {
    return (
      <div className="projectMembers">
        <p></p>
        <h1>Members</h1>
        <p>Cupcake ipsum dolor sit amet chocolate bar sesame snaps sugar plum dessert. Sugar plum sesame snaps oat cake jelly cake sugar plum cake danish pie. Jelly-o candy canes soufflé gummi bears jelly beans sweet roll bear claw.</p>
        <Table>
          <Row header={true}>
            <Cell>USER</Cell>
            <Cell>ROLE</Cell>
            <Cell>ACCESS</Cell>
          </Row>
          <Row>
            <Cell>First Last</Cell>
            <Cell>Admin</Cell>
            <Cell>Owner</Cell>
          </Row>
          <Row>
            <Cell>First Last</Cell>
            <Cell>Project Lead</Cell>
            <Cell>Collaborator</Cell>
          </Row>
        </Table>
      </div>
    );
  }
}

function mapStateToProps(state) {
  return {
    members: state.collaborators.members,
  };
}

function mapDispatchToProps(dispatch) {
  return {
    fetchMembers: bindActionCreators(fetchMembers, dispatch)
  };
}

export default connect<any, DispatchProps, any>(mapStateToProps, mapDispatchToProps)(ProjectMembers);
