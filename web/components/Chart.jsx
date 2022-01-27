import React from "react"
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
    const labels = ["GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"]

    const data = {
        labels,
        datasets: [
            {
                label: "Methods",
                data: [1, 2, 3, 1, 10],
                backgroundColor: "rgba(255, 99, 132, 0.5)",
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

    return <Bar options={options} data={data} />
}
