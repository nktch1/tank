import http from "k6/http";

export default function() {
    let response = http.get("https://test-api.k6.io");
};

export let options = {
    vus: 10,
    stages: [
        { duration: "1s", target: 100 },
    ]
};
