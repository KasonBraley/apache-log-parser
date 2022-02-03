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

import { getChartData } from "./data"

ChartJS.register(CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend)

export default function BarChart({ type, data }) {
    let chartData = getChartData(type, data)

    return <Bar options={chartData.options} data={chartData.chartData} />
}
