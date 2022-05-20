import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
    stages: [
        { duration: '2.5m', target: 1000 },
        { duration: '5m', target: 1000 },
        { duration: '2.5m', target: 0 },
    ],
    thresholds: {
        'http_req_duration': ['p(99)  < 1000'],
    },
    summaryTrendStats: ['min', 'avg', 'p(99)', 'max'],
};

const BASE_URL = 'http://localhost:8080';

export default () => {
    const res = http.get(`${BASE_URL}/apps`).json();

    check(res, { 'retrieved games': (obj) => obj.data.apps.length > 0 });

    sleep(1);
};
