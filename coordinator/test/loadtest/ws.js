import ws from 'k6/ws';
import { check } from 'k6';

export const options = {
    stages: [
        { duration: '3m', target: 1000 },
        { duration: '4m', target: 1000 },
        { duration: '3m', target: 0 },
    ],
    summaryTrendStats: ['min', 'avg', 'p(99)', 'max'],
};

const getRandomInt = (max) => {
    return Math.floor(Math.random() * max);
}

const randUser = () => {
    const randInt = getRandomInt(2);

    if (randInt === 0) {
        return {
            type: 'join',
            data: JSON.stringify({
                role: "player",
            }),
        }
    }

    return {
        type: 'join',
        data: JSON.stringify({
            role: "provider",
        }),
    }
}

export default () => {
    const url = 'ws://localhost:8080/ws';

    const res = ws.connect(url, null, function (socket) {
        socket.on('open', () => socket.send(JSON.stringify(randUser())));
    });

    check(res, { 'status is 101': (r) => r && r.status === 101 });
}
