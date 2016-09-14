/**
 * Created by justin on 6/25/16.
 */

import * as React from 'react';
import * as $ from 'jquery';
import * as classNames from 'classnames';
import { Link, withRouter } from 'react-router';
import { Sidebar } from '../Sidebar/Sidebar';
import { buildPath } from '../../../App/utils/buildPath';
import { getRoute } from '../../../App/utils/getRoute';
import { routes } from '../../../routes';
import * as _ from 'lodash';
import { connect } from 'react-redux';
import { fetchProfile } from '../../../Profile/actions/profile.actions';
import { bindActionCreators } from 'redux';
import './navigation.scss';
import { Project } from '../../../Proxy/Proxy';
const logo = require('../../../../assets/h2o-home.png');
import {Motion, spring} from 'react-motion';

interface Props {
  routes: any
  params: any
  profile: {
    isEulaAgreed: boolean
  },
  project: Project
}

interface DispatchProps {
  fetchProfile: Function
}


interface State {
  activeTopLevelPath: string
  isSubMenuActive: boolean
  isEulaAgreed: boolean
}

export class Navigation extends React.Component<Props & DispatchProps, any> {

  constructor() {
    super();
    this.state = {
      activeTopLevelPath: '',
      isSubMenuActive: false,
      isEulaAgreed: false
    };
  }

  componentWillMount(): void {
    this.setMenuState(this.props.routes);
    this.props.fetchProfile();
  }

  componentWillReceiveProps(nextProps: Props): void {
    this.setMenuState(nextProps.routes);
    this.setState({
      isEulaAgreed: nextProps.profile.isEulaAgreed as boolean
    });
  }

  setMenuState(newRoutes: any[]): void {
    let currentRoutePath = newRoutes[newRoutes.length - 1].path;
    let topLevelPath = '';
    if (currentRoutePath) {
      topLevelPath = currentRoutePath.split('/')[0];
    }
    let submenuActive = false;
    _.forEach(routes[0].childRoutes, (route) => {
      if (this.isActive(route.path, newRoutes) && route.showChildrenAsSubmenu) {
        submenuActive = true;
      }
    });

    this.setState({
      activeTopLevelPath: topLevelPath,
      isSubMenuActive: submenuActive
    });
  }

  isActive(path: string, newRoutes: any[]): boolean {
    let currentRoutePath = newRoutes[newRoutes.length - 1].path;
    if (currentRoutePath && currentRoutePath.indexOf(path) !== -1) {
      return true;
    }
    return false;
  }

  getParentRouteName(currentPath: string): string {
    let newPathParts = currentPath.split('/');
    if (newPathParts.length < 2) {
      return currentPath;
    }
    newPathParts.pop();
    let newPath = newPathParts.join('/');
    let parentRoute = getRoute(newPath);
    return parentRoute.name;
  }

  logout() {
    $.ajax({
      url: window.location.protocol + '://' + window.location.host,
      beforeSend: function (xhr) {
        xhr.withCredentials = true;
        xhr.setRequestHeader('Authorization', 'Basic ' + btoa('fjkdshfhkjsdfjkhsdkfjhsdf:hfkjdshfdhff'));
      }
    });
  }

  renderSubmenu(activeRoute: any, shouldShow: boolean): JSX.Element {
    let childRoutes = routes[0].childRoutes.filter((route) => {
      return (route.path.indexOf(activeRoute.path) !== -1 && route.path !== activeRoute.path);
    });
    let styleTo;
    console.log(shouldShow);
    if(shouldShow) {
      styleTo = { left: spring(72) }
    } else {
      styleTo = { left: spring(300) }
    }
    return (
    <Motion defaultStyle={{left: 300}} style={styleTo}>
      {interpolatingStyle =>
      <div className="left-submenu" style={interpolatingStyle}>
        <Sidebar className='secondary-navigation'>
          <nav className="navigation--primary">
            <div className="navigation">
              <header>
                <div className="header-navigation">
                  <Link to={this.getParentRouteName(activeRoute.path)}><i
                    className="fa fa-angle-left"></i><span>{this.getParentRouteName(activeRoute.path)}</span></Link>
                </div>
              </header>
              <div className="header-content">{this.props.project.name}</div>
              <ul className="nav-list">
                {_.map(childRoutes, (menuItem: any) => {
                  let path = buildPath(menuItem.path, this.props.params);
                  return (!menuItem.showInNavigation) ? null : (
                    <li key={menuItem.path}
                        className={classNames('nav-list--item', {active: this.isActive(menuItem.path, this.props.routes)})}>
                      <Link to={path}>{menuItem.name}</Link>
                    </li>
                  );
                })}
              </ul>
            </div>
          </nav>
        </Sidebar>
      </div>
      }
    </Motion>
    );
  }

  render(): React.ReactElement<HTMLElement> {
    let submenu = <div></div>;
    return (
      <div className={classNames('nav-container', {hidden: !this.state.isEulaAgreed})}>
        <Sidebar className="primary-navigation">
          <nav className="navigation--primary">
            <div className="navigation">
              <header>
                <div className="logo-container">
                  <Link to="/">
                    <div className="logo">STEAM</div>
                  </Link>
                </div>
              </header>
              <div className="header-content">
              </div>
              <ul className='nav-list'>
                {routes[0].childRoutes.map((route: any) => {
                  let isActive = false;
                  if (this.isActive(route.path, this.props.routes)) {
                    isActive = true;
                    if (route.showChildrenAsSubmenu) {
                      submenu = this.renderSubmenu(route, isActive);
                    }
                  }
                  if (route.path.split('/').length > 1 || !route.showInNavigation) {
                    return null;
                  }
                  let activeChildren = route.path === this.state.activeTopLevelPath && this.state.isSubMenuActive;
                  let path = '/' + route.path;
                  return (
                    <li key={path}
                        className={classNames('nav-list--item', { active: isActive}, {activeChildren: activeChildren}) }>
                      <Link to={path}><i className={route.icon}></i>
                        <div className="nav-list--label">{route.name}</div>
                      </Link>
                    </li>
                  );
                })
                }
                <li className="logout nav-list--item">
                  <a href="mailto:steam@h2o.ai?subject=STEAM: ">
                    <i className="fa fa-question-circle-o"/>
                    <div className="nav-list--label">
                      Support
                    </div>
                  </a>
                  <a onClick={this.logout.bind(this)}>
                    <i className="fa fa-sign-out"/>
                    <div className="nav-list--label">Logout</div>
                  </a>
                </li>
              </ul>
            </div>
          </nav>
        </Sidebar>
        {submenu}
      </div>
    );
  }
}


function mapStateToProps(state): any {
  return {
    project: state.projects.project,
    profile: state.profile
  };
}

function mapDispatchToProps(dispatch) {
  return {
    fetchProfile: bindActionCreators(fetchProfile, dispatch)
  };
}

export default connect<any, DispatchProps, any>(mapStateToProps, mapDispatchToProps)(withRouter(Navigation));
