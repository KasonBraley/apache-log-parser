import Head from "next/head"
import { useState } from "react"
import BarChart from "../components/Chart"
import Table from "../components/Table"

export default function Home() {
    const [file, setFile] = useState()
    const [view, setView] = useState("chart")

    function handleChange(event) {
        setFile(event.target.files[0])
    }

    function handleSubmit(e) {
        e.preventDefault()
        sendData("http://localhost:5000/upload", file)
    }

    async function sendData(url, data) {
        const formData = new FormData()

        formData.append("file", data)

        const response = await fetch(url, {
            method: "POST",
            body: formData,
        })

        console.log(await response.json())
    }

    return (
        <div className="">
            <Head>
                <title>Apache Log Parser</title>
                <meta name="Apache Log Parser and Aggregator" />
            </Head>

            <main className="flex justify-center items-center min-h-screen space-x-60">
                {/* <span>Example Apache Common log line</span> */}
                {/* <p> */}
                {/*     { */}
                {/*         '132.128.161.195 - - [25/Jan/2022:20:08:53 -0700] "HEAD /synergize/deploy/cutting-edge/convergence HTTP/2.0" 301 14575' */}
                {/*     } */}
                {/* </p> */}

                <form
                    onSubmit={handleSubmit}
                    className="border-2 border-black border-solid p-4"
                >
                    <h1 className="text-lg mb-4">Apache Log Upload</h1>
                    <label className="block">
                        <span className="sr-only">Choose Apache log file</span>
                        <input
                            type="file"
                            onChange={handleChange}
                            className="block w-full text-sm text-slate-500
        file:mr-4 file:py-2 file:px-4
        file:rounded-full file:border-0
        file:text-sm file:font-semibold
        file:bg-violet-50 file:text-violet-700
        hover:file:bg-violet-100
      "
                        />
                    </label>

                    <button type="submit" className="btn btn-primary">
                        Upload
                    </button>
                </form>

                <div>
                    <div className="navbar mb-2 shadow-lg bg-neutral text-neutral-content rounded-box">
                        <div className="flex-1 px-2 mx-2">
                            <a
                                onClick={() => setView("chart")}
                                className="btn btn-ghost btn-sm rounded-btn"
                            >
                                Chart
                            </a>
                        </div>
                        <div className="flex-1 px-2 mx-2">
                            <a
                                onClick={() => setView("table")}
                                className="btn btn-ghost btn-sm rounded-btn"
                            >
                                Table
                            </a>
                        </div>
                    </div>

                    {view === "chart" ? <BarChart /> : <Table />}
                </div>
            </main>
        </div>
    )
}
