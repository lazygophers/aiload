import {message} from 'antd';
import axios from 'axios';
import cookie from 'react-cookies';

const {getLanguage} = await import('@i18n/tran');

axios.defaults.retry = 3;
axios.defaults.retryDelay = 1000 * 60 * 5;

export function onLogout() {
	cookie.remove('token');
	cookie.remove('user');
	window.location.href = '';
}

export let service = axios.create({
	// baseURL: location.origin,
	timeout: 60000,
	maxRate: 1,
	maxRedirects: 10,
	withCredentials: true
});

service.interceptors.request.use(
	(config) => {
		console.debug(
			`request ${config.method} ${config.url} ${JSON.stringify(config.data)}`
		);
		config.headers['X-Token'] = cookie.load('token');
		config.headers['X-Client'] = 'website';
		config.headers['X-Language'] = getLanguage();
		return config;
	},
	(error) => {
		message.error(error);
		console.log(error); // for debug
		return Promise.reject(error);
	}
);

// response interceptor
service.interceptors.response.use(
	(response) => {
		if (response.status === 200) {
			const res = response.data; //res is my own data
			if (!res.code || res.code === 0) {
				console.debug(
					`response ${response.config.method} ${response.config.url} ${JSON.stringify(res.data)}`
				);
				return res.data;
			} else if (res.code == 401 || res.code == 403 || res.code == 1002) {
				console.log('logout because of code');
				onLogout();
				message.error(res.message || 'error');
				return Promise.reject(
					new Error(res.error || res.message || 'Error')
				);
			} else {
				console.debug(
					`response ${response.config.method} ${response.config.url} ${res}`
				);
				message.error(res.message || 'error');
				return Promise.reject(
					new Error(res.error || res.message || 'Error')
				);
			}
		} else if (response.status == 401 || response.status == 403) {
			console.log('logout because of status');
			onLogout();
			return Promise.reject(new Error(response.status));
		} else {
			message.error(response.status);
			return Promise.reject(new Error(response.status));
		}
	},
	(error) => {
		if (error.response.status == 401) {
			console.log('logout because of error status');
			onLogout();
			return Promise.reject(new Error(error.response.status));
		}
		message.error(error.response.status);
		return Promise.reject(error);
	}
);

export default service; //导出封装后的axios

export function getBaseUrl() {
	return window.location.origin;
}

export function getWsBaseUrl() {
	if (window.location.protocol === 'https') {
		return (
			'wss://' +
			window.location.hostname +
			':' +
			window.location.port +
			'/ws'
		);
	}

	return (
		'ws://' + window.location.hostname + ':' + window.location.port + '/ws'
	);
}

export async function post(path, body) {
	return service({
		url: getBaseUrl() + '/api' + path,
		method: 'POST',
		data: body,
		headers: {
			'Content-Type': 'application/json; charset=utf-8'
		}
	});
}

export async function put(path, body) {
	return service({
		url: getBaseUrl() + '/api' + path,
		method: 'PUT',
		data: body,
		headers: {
			'Content-Type': 'application/json; charset=utf-8'
		}
	});
}

export async function get(path, params) {
	return service({
		url: getBaseUrl() + '/api' + path,
		method: 'GET',
		params: params
	});
}

export async function del(path, params) {
	return service({
		url: getBaseUrl() + '/api' + path,
		method: 'DELETE',
		params: params
	});
}
