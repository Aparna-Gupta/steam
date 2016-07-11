/**
 * Created by justin on 7/5/16.
 */

import * as React from 'react';
import '../styles/breadcrumb.scss';

interface Props {
  crumbs: string[]
}

export default class Breadcrumb extends React.Component<Props, any> {
  isHiddenBreadcrumb() {
    let result = _.some(this.props.routes, route => {
      return route.isHiddenBreadcrumb === true
    });
    console.log(result);
    return result;
  }

  render(): React.ReactElement<HTMLElement> {
    if (this.isHiddenBreadcrumb() === true) {
      return null;
    }
    return (
      <ol className="breadcrumb">
        {this.props.routes.map((route, i) => {
          console.log(route);
          if (route.isExcludedFromBreadcrumb === true) {
            return null;
          }
          return <li key={i}>{route.name}</li>;
        })}
      </ol>
    );
  }
}
