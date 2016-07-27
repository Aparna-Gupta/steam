/**
 * Created by justin on 7/10/16.
 */
import * as React from 'react';
import * as classNames from 'classnames';
import DefaultModal from '../../App/components/DefaultModal';
import { Link } from 'react-router';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { fetchProfile, setProfile } from '../../Profile/actions/profile.actions';
import '../styles/welcomesplashscreen.scss';
import { EULA_TEXT } from '../utils/eula.text';

interface Props {
  profile: any
}

interface DispatchProps {
  fetchProfile: Function,
  setProfile: Function
}

export class WelcomeSplashScreen extends React.Component<Props & DispatchProps, any> {
  constructor() {
    super();
    this.state = {
      isEulaAgreed: false
    };
  }

  componentWillReceiveProps(nextProps) {
    this.setState({
      isEulaAgreed: nextProps.profile.isEulaAgreed
    });
  }

  onChangeHandler(event) {
    let agreed = event.target.checked;
    this.setState({
      isEulaAgreed: agreed
    }, () => {
      this.props.setProfile({
        isEulaAgreed: agreed
      });
    });
  }

  openEula() {
    this.setState({
      isEulaOpen: true
    });
  }

  agree() {
    this.setState({
      isEulaAgreed: true
    }, () => {
      this.props.setProfile({
        isEulaAgreed: true
      });
    });
    this.close();
  }

  close() {
    this.setState({
      isEulaOpen: false
    });
  }

  render(): React.ReactElement<HTMLDivElement> {
    console.log(this.props.profile, this.state.isEulaAgreed);
    return (
      <div className="welcome-splash-screen">
        <DefaultModal className="eula-modal" open={this.state.isEulaOpen}>
          <div className="eula-content">
            <div className="eula-text">
              {EULA_TEXT}
            </div>
            <div className="buttons-container">
              <button className="default" onClick={this.agree.bind(this)}>I Agree</button>
              <button className="default" onClick={this.close.bind(this)}>Cancel</button>
            </div>
          </div>
        </DefaultModal>
        <div className="welcome-splash-screen--content">
          <div>WELCOME TO</div>
          <div className="welcome-splash-screen--content--product-name">H<sub>2</sub>O STEAM</div>
          <div>Fast, distributed data science for teams</div>
          <div><Link to="/projects/new"
                     className={classNames('default', 'start-project', {disabled: !this.state.isEulaAgreed})}>Start A New Project</Link>
          </div>
          <div className="eula-line">
            <input type="checkbox" checked={this.state.isEulaAgreed} onChange={this.onChangeHandler.bind(this)}/>
            I understand this is a preview release, and agree with all <a href="javascript:void(0);"
                                                                          onClick={this.openEula.bind(this)}>terms and conditions</a>.
          </div>
        </div>
      </div>
    );
  }
}

function mapStateToProps(state): any {
  return {
    profile: state.profile
  };
}

function mapDispatchToProps(dispatch): DispatchProps {
  return {
    fetchProfile: bindActionCreators(fetchProfile, dispatch),
    setProfile: bindActionCreators(setProfile, dispatch)
  };
}

export default connect<any, DispatchProps, any>(mapStateToProps, mapDispatchToProps)(WelcomeSplashScreen);
