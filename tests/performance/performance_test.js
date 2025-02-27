import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  scenarios: {
    burst_submit_questionnaire: {
      executor: 'constant-vus',
      vus: 200,
      duration: '100s',
      exec: 'burst_submit_questionnaire',
      tags: { scenario: 'burst_submit_questionnaire' },
      startTime: '0s',
    },
    burst_get_active_questionnaire: {
      executor: 'ramping-vus',
      startVUs: 10,
      exec: 'burst_get_active_questionnaire',
      tags: { scenario: 'burst_get_active_questionnaire' },
      stages: [
        { target: 350, duration: '100s' },
      ],
      gracefulRampDown: '30s',
      startTime: '101s',
    },
  },
  thresholds: {
    'http_req_duration{scenario:burst_submit_questionnaire}': [
        'p(95)<2000',
        'p(99)<3000',
        'avg<200',
    ],
    'http_req_failed{scenario:burst_submit_questionnaire}': [
        'rate==0',
    ],

    'http_req_duration{scenario:burst_get_active_questionnaire}': [
        'p(95)<200',
        'p(99)<300',
        'avg<100',
    ],
    'http_req_failed{scenario:burst_get_active_questionnaire}': ['rate==0'],
  },
};

// fetch the available active questionnaire to be used to submit
export function setup() {
    const url = "http://atec_api:5000/v1/atec/packages/active";
    const headers = {
        "accept": "application/json",
    };

    const response = http.get(url, { headers: headers });
    if (response.status !== 200) {
        throw new Error("failed to fetch active questionnaire");
    }

    const body = JSON.parse(response.body);
    const data = body.data;

    if (data.length === 0) {
        throw new Error("no active questionnaire found");
    }

    const activeQuestionnaire = data[0];

    return {
        'activeQuestionnaireId': activeQuestionnaire.id,
    };
}

// The default exported function is gonna be picked up by k6 as the entry point for the test script. It will be executed repeatedly in "iterations" for the whole duration of the test.
export function burst_submit_questionnaire(data) {
    const { activeQuestionnaireId } = data;
    const url = "http://atec_api:5000/v1/atec/questionnaires"
    const payload = {
        answers: {
            0: {
                1: 0,
                2: 0,
                3: 0,
                4: 0,
                5: 0,
                6: 0,
                7: 0,
                8: 0,
                9: 0,
                10: 0,
                11: 0,
                12: 0,
                13: 0,
                14: 0,
            },
            1: {
                1: 0,
                2: 0,
                3: 0,
                4: 0,
                5: 0,
                6: 0,
                7: 0,
                8: 0,
                9: 0,
                10: 0,
                11: 0,
                12: 0,
                13: 0,
                14: 0,
                15: 0,
                16: 0,
                17: 0,
                18: 0,
                19: 0,
                20: 0,
            },
            2: {
                1: 0,
                2: 0,
                3: 0,
                4: 0,
                5: 0,
                6: 0,
                7: 0,
                8: 0,
                9: 0,
                10: 0,
                11: 0,
                12: 0,
                13: 0,
                14: 0,
                15: 0,
                16: 0,
                17: 0,
                18: 0,
            },
            3: {
                1: 0,
                2: 0,
                3: 0,
                4: 0,
                5: 0,
                6: 0,
                7: 0,
                8: 0,
                9: 0,
                10: 0,
                11: 0,
                12: 0,
                13: 0,
                14: 0,
                15: 0,
                16: 0,
                17: 0,
                18: 0,
                19: 0,
                20: 0,
                21: 0,
                22: 0,
                23: 0,
                24: 0,
                25: 0,
            },
        },
        package_id: activeQuestionnaireId,
    };

    const res = http.post(url, JSON.stringify(payload),
        {
            headers: {
                "accept": "application/json",
                "Content-Type": "application/json",
            },
        },
        {
            tags: { scenario: 'burst_submit_questionnaire' }
        }
    );

    check(res, {
        "burst_submit_questionnaire response must 200": (res) => res.status == 200,
    });
}

export function burst_get_active_questionnaire () {
  const res = http.get('http://atec_api:5000/v1/atec/packages/active',
        {
            tags: { scenario: 'burst_get_active_questionnaire' },
        },
    );

  check(res, {
        "burst_get_active_questionnaire response must 200": (res) => res.status == 200,
    });
}