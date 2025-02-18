import http from 'k6/http';

export const options = {
  stages: [
    { duration: '10s', target: 5},
    { duration: '30s', target: 10 },
    { duration: '10s', target: 50 },
    { duration: '30s', target: 14 },
    { duration: '30s', target: 10 },
    { duration: '30s', target: 10 },
    { duration: '2s', target: 100 },
    { duration: '3s', target: 200 },
    { duration: '3s', target: 300 },
    { duration: '4s', target: 400 },
    { duration: '5s', target: 500 },
    { duration: '5s', target: 600 },
    { duration: '5s', target: 700 },
    { duration: '5s', target: 800 },
    { duration: '5s', target: 900 },
    { duration: '5s', target: 1000 },
  ],
};

// The default exported function is gonna be picked up by k6 as the entry point for the test script. It will be executed repeatedly in "iterations" for the whole duration of the test.
export default function () {
  // Make a GET request to the target URL
  http.get('https://atec.luckyakbar.web.id/v1/atec/packages/active');
}