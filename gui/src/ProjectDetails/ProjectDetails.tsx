/**
 * Created by justin on 6/27/16.
 */

import * as React from 'react';
import * as classNames from 'classnames';
import Collapsible from './components/Collapsible';
import ModelOverview from './components/ModelOverview';
import GoodnessOfFit from './components/GoodnessOfFit';
import './styles/projectdetails.scss';
import 'bootstrap/dist/css/bootstrap.css';

interface Props {
  params: any
}


export default class ProjectDetails extends React.Component<Props, any> {
  constructor() {
    super();
    this.state = {
      isModelOpen: true,
      isResidualOpen: true,
      isVariableOpen: true,
      isGoodnessOpen: true
    };
  }

  toggleOpen(accordian: string) {
    /**
     * TODO(justinloyola): Fix the asynchronous state change issues
     */
    if (accordian === 'model') {
      this.setState({
        isModelOpen: !this.state.isModelOpen
      });
    } else if (accordian === 'residual') {
      this.setState({
        isResidualOpen: !this.state.isResidualOpen
      });
    } else if (accordian === 'variable') {
      this.setState({
        isVariableOpen: !this.state.isVariableOpen
      });
    } else if (accordian === 'goodness') {
      this.setState({
        isGoodnessOpen: !this.state.isGoodnessOpen
      });
    }
  }

  render(): React.ReactElement<HTMLDivElement> {
    return (
      <div className="project-details">
        <div className="model-header">
          Model: GBT-1069085
        </div>
        <header>
          <span onClick={this.toggleOpen.bind(this, 'model')}><i
            className={classNames('fa', {'fa-minus-square-o': this.state.isModelOpen, 'fa-plus-square-o': !this.state.isModelOpen})}></i
          >Model Overview</span>
        </header>
        <Collapsible open={this.state.isModelOpen}>
          <ModelOverview></ModelOverview>
        </Collapsible>
        <header>
          <span onClick={this.toggleOpen.bind(this, 'goodness')}><i
            className={classNames('fa', {'fa-minus-square-o': this.state.isGoodnessOpen, 'fa-plus-square-o': !this.state.isGoodnessOpen})}></i
          >Goodness of Fit</span>
        </header>
        <Collapsible open={this.state.isGoodnessOpen}>
          <GoodnessOfFit></GoodnessOfFit>
        </Collapsible>
        <header>
          <span onClick={this.toggleOpen.bind(this, 'residual')}><i
            className={classNames('fa', {'fa-minus-square-o': this.state.isResidualOpen, 'fa-plus-square-o': !this.state.isResidualOpen})}></i
          >Residual Analysis</span>
        </header>
        <Collapsible open={this.state.isResidualOpen}>
          <div>
            Residual body
          </div>
        </Collapsible>
        <header>
          <span onClick={this.toggleOpen.bind(this, 'variable')}><i
            className={classNames('fa', {'fa-minus-square-o': this.state.isVariableOpen, 'fa-plus-square-o': !this.state.isVariableOpen})}></i
          >Variable Analysis</span>
        </header>
        <Collapsible open={this.state.isVariableOpen}>
          <div>
            Variable body
          </div>
        </Collapsible>
      </div>
    );
  }
}