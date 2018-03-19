import React, { Component } from 'react';

import {
	Sparkline,
	LineSeries,
	HorizontalReferenceLine,
	BandLine,
	PatternLines,
	PointSeries } from '@data-ui/sparkline';
import { allColors } from '@data-ui/theme'; // open-color colors

const data = Array(25).fill().map(Math.random);


var foo = 0;
class Chart extends Component {
	constructor(props) {
		super(props);
		this.state = {
			data: Array(25).fill().map(Math.random)
		};
	}

	componentDidMount() {
		setInterval(() => {
			var d = this.state.data.slice(1, this.state.data.size);
			var x = Math.random();

			d.push(x);

			this.setState({data: d});
		}, 200);
	}

	render() {
		return (
			<div>
			<Sparkline
			ariaLabel="A line graph of randomly-generated data"
			margin={{ top: 24, right: 64, bottom: 24, left: 64,}}
			width={800}
			height={100}
			data={this.state.data}
			valueAccessor={datum => datum}
			>
			{/* this creates a <defs> referenced for fill */}
			<PatternLines
			id="unique_pattern_id"
			height={6}
			width={6}
			stroke="#000"
			strokeWidth={1}
			orientation={['diagonal']}
			/>
			{/* display innerquartiles of the data */}
			<BandLine
			band="innerquartiles"
			fill="url(#unique_pattern_id)"
			/>
			{/* display the median */}
			<HorizontalReferenceLine
			stroke="#000"
			strokeWidth={1}
			strokeDasharray="4 4"
			reference="median"
			/>
			{/* Series children are passed the data from the parent Sparkline */}
			<LineSeries
			showArea={false}
			stroke="#000"
			/>
			<PointSeries
			points={['min', 'max']}
			fill="black"
			size={5}
			stroke="#fff"
			renderLabel={val => val.toFixed(2)}
			/>
			</Sparkline>
			</div>
		);
	}
}

export default Chart;
