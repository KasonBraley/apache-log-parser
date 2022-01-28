import React, { useEffect, useState } from "react"
import { Bar } from "react-chartjs-2"
import {
    Chart as ChartJS,
    CategoryScale,
    LinearScale,
    BarElement,
    Title,
    Tooltip,
    Legend,
} from "chart.js"

ChartJS.register(CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend)

export default function BarChart() {
    const labels = ["DELETE", "GET", "HEAD", "PATCH", "POST", "PUT"]
    let [data, setData] = useState()

    useEffect(() => {
        async function getData() {
            const response = await fetch("http://localhost:4003/retrieve")
            if (response.ok) {
                let resp = await response.json()
                console.log(resp)
                let methods = resp.map((log) => log.Method)

                var map = methods.reduce(function (prev, cur) {
                    prev[cur] = (prev[cur] || 0) + 1
                    return prev
                }, {})

                setData(map)
            } else {
                console.log("ERROR fetching the database data")
            }
        }

        getData()
    }, [])

    const chartData = {
        labels,
        datasets: [
            {
                label: "Methods",
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
        ],
    }

    const options = {
        responsive: true,
        plugins: {
            legend: {
                position: "top",
            },
            title: {
                display: true,
                text: "# of HTTP Methods",
            },
        },
    }

    return <Bar options={options} data={chartData} />
}
