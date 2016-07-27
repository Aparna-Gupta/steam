/**
 * Created by justin on 6/25/16.
 */

import * as React from 'react';
import * as classNames from 'classnames';
import { Link, withRouter } from 'react-router';
import { Sidebar } from '../Sidebar/Sidebar';
import './navigation.scss';
import { routes } from '../../../routes';
import * as _ from 'lodash';
import { connect } from 'react-redux';
const logo = require('../../../../assets/h2o-home.png');

interface Props {
  router: any,
  routes: any
}

interface DispatchProps {
}

export class Navigation extends React.Component<Props & DispatchProps, any> {

  constructor() {
    super();
    this.state = {
      isHidden: localStorage.getItem('steamDidAgreeToEula') === 'false',
      isEulaAgreed: false
    };
  }

  sitemap = routes[0].childRoutes;

  componentWillReceiveProps(nextProps) {
    this.setState({
      isEulaAgreed: nextProps.profile.isEulaAgreed
    });
  }

  render(): React.ReactElement<HTMLElement> {

    let submenu = null;

    _.forEach(this.props.routes, (route) => {
      if (!submenu && route.path && route.path !== "/" && this.props.router.isActive(route.path)) {
        submenu = (
          <Sidebar className='secondary-navigation'>
            <nav className="navigation--primary">
              <div className="navigation">
                <header>
                  <div className="header-navigation">
                    <i className="fa fa-angle-left"></i><span>{route.name}</span>
                  </div>
                </header>
                <div className="header-content">UNTITLED</div>
                <ul className="nav-list">
                  {_.map(route.childRoutes, (menuItem: any) => {
                    let path = route.path + '/' + menuItem.path;
                    return (!menuItem.showInNavigation) ? null : (
                      <li key={menuItem.path} className={classNames('nav-list--item', {active: this.props.router.isActive(path)})}>
                        <Link to={path}>{menuItem.name}</Link>
                      </li>
                    );
                  })}
                </ul>
              </div>
            </nav>
          </Sidebar>
        );
      }
    });

    return (
      <div className={classNames('nav-container', {hidden: !this.state.isEulaAgreed})}>
        <Sidebar className="primary-navigation">
          <nav className="navigation--primary">
            <div className="navigation">
              <header>
                <div className="logo-container">
                  <Link to="/"><div className="logo"><img src={logo}></img></div></Link>
                </div>
              </header>
              <div className="header-content">
              </div>
              <ul className='nav-list'>
              {_.map(this.sitemap, (route: any) => {
                  return (
                    <li key={route.path} className={classNames('nav-list--item', {active: this.props.router.isActive(route.path)})}>
                      <Link to={route.path}><i className={route.icon}></i><div className="nav-list--label">{route.name}</div></Link>
                    </li>
                  );
                })
              }
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
    profile: state.profile
  };
}

export default connect<any, DispatchProps, any>(mapStateToProps, {})(withRouter(Navigation));
