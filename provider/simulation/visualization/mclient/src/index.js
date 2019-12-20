import { render } from 'react-dom';
import { getCombinedReducer } from 'harmoware-vis';
import { createStore } from 'redux';

import { Provider } from 'react-redux';
import React from 'react';
import App from './containers/app';
import 'bootstrap/dist/css/bootstrap.min.css';
import './scss/harmovis.scss';

const store = createStore(getCombinedReducer());

render(
	<Provider store={store}>
	<App />
	</Provider>,
    document.getElementById('app')
);
