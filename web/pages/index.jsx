import Head from "next/head"
import React, { useEffect, useState } from "react"
import BarChart from "../components/Chart/Chart"
import Table from "../components/Table"

export default function Home() {
    const [file, setFile] = useState()
    const [view, setView] = useState("chart")
    const logInputRef = React.useRef()

    let [methods, setMethods] = useState()
    let [statusCodes, setStatusCodes] = useState()
    let [httpVersions, setHttpVersions] = useState()

    useEffect(() => {
        getData()
    }, [])

    function handleChange(event) {
        setFile(event.target.files[0])
    }

    function handleSubmit(e) {
        e.preventDefault()
        sendData("http://localhost:4000/upload", file)
        logInputRef.current.value = "" //Resets the file name of the file input
    }

    async function getData() {
        const response = await fetch("http://localhost:4001/retrieve")
        if (response.ok) {
            let resp = await response.json()
            console.log(resp)
            let methods = []
            let statusCodes = []
            let httpVersion = []

            resp.map((log) => {
                methods.push(log.Method)
                statusCodes.push(log.Status)
                httpVersion.push(log.HTTPVersion)
            })

            let reducedMethods = reduceData(methods)
            let reducedStatusCodes = reduceData(statusCodes)
            let reducedHttpVersions = reduceData(httpVersion)

            setMethods(reducedMethods)
            setStatusCodes(reducedStatusCodes)
            setHttpVersions(reducedHttpVersions)
        } else {
            console.log("ERROR fetching the database data")
        }
    }

    function reduceData(arr) {
        return arr.reduce(function (prev, cur) {
            prev[cur] = (prev[cur] || 0) + 1
            return prev
        }, {})
    }

    async function sendData(url, data) {
        const formData = new FormData()

        formData.append("file", data)

        const response = await fetch(url, {
            method: "POST",
            body: formData,
        })

        if (response.ok) {
            console.log(await response.text())
            getData()
        }
    }

    return (
        <>
            <Head>
                <title>Apache Log Parser</title>
                <meta name="Apache Log Parser and Aggregator" />
            </Head>

            <main className="flex flex-col 2xl:flex 2xl:flex-row justify-center items-center min-h-screen min-w-full space-y-40 2xl:space-y-0 2xl:space-x-52">
                <form
                    onSubmit={handleSubmit}
                    className="w-64 2xl:w-96 border-2 border-black border-solid p-4 rounded-md shadow-xl shadow-violet-500/50"
                >
                    <h1 className="text-lg mb-4">Apache Log Upload</h1>
                    <label className="block">
                        <span className="sr-only">Choose Apache log file</span>
                        <input
                            type="file"
                            onChange={handleChange}
                            ref={logInputRef}
                            accept=".log, .txt"
                            className="block w-full text-sm text-slate-200
                                file:mr-4 file:py-2 file:px-4
                                file:rounded-full file:border-0
                                file:text-sm file:font-semibold
                                file:bg-violet-50 file:text-violet-700
                                hover:file:bg-violet-100"
                        />
                    </label>

                    <button
                        type="submit"
                        className="btn btn-primary mt-6 rounded-full"
                    >
                        Upload
                    </button>
                </form>

                <div className="w-36 tabs tabs-boxed">
                    <a
                        onClick={() => setView("chart")}
                        className={view === "chart" ? "tab tab-active" : "tab"}
                    >
                        Chart
                    </a>
                    <a
                        onClick={() => setView("table")}
                        className={view === "table" ? "tab tab-active" : "tab"}
                    >
                        Table
                    </a>
                </div>

                <div className="w-[600px]">
                    {view === "chart" ? (
                        <>
                            <BarChart type="methods" data={methods} />
                            <BarChart type="statusCodes" data={statusCodes} />
                            <BarChart type="httpVersions" data={httpVersions} />
                        </>
                    ) : (
                        <Table />
                    )}
                </div>
            </main>
        </>
    )
}
