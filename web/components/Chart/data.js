export function getChartData(type, data) {
    let chartData = setChartData(type, data, labels[type], datasets(type, data))
    let options = setOptions(type)

    return {
        chartData,
        options,
    }
}

const labels = {
    methods: ["DELETE", "GET", "HEAD", "PATCH", "POST", "PUT"],
    statusCodes: [200, 201, 204, 301, 304, 406, 502],
    httpVersions: [1, 2],
}

function datasets(type, data) {
    switch (type) {
        case "methods":
            return [
                {
                    label: type,
                    data: data && [
                        data.DELETE,
                        data.GET,
                        data.HEAD,
                        data.PATCH,
                        data.POST,
                        data.PUT,
                    ],
                    backgroundColor: "rgba(122, 75, 208, 1)",
                },
            ]
        case "statusCodes":
            return [
                {
                    label: type,
                    data: data && [
                        data[200],
                        data[201],
                        data[204],
                        data[301],
                        data[304],
                        data[406],
                        data[502],
                    ],
                    backgroundColor: "rgba(122, 75, 208, 1)",
                },
            ]
        case "httpVersions":
            return [
                {
                    label: type,
                    data: data && [data[1], data[2]],
                    backgroundColor: "rgba(122, 75, 208, 1)",
                },
            ]
        default:
            return []
    }
}

function setChartData(type, data, labels, datasets) {
    return {
        labels,
        datasets,
    }
}

function setOptions(type) {
    return {
        responsive: true,
        plugins: {
            legend: {
                position: "top",
            },
            title: {
                display: true,
                text: `# of ${type}`,
            },
        },
    }
}
