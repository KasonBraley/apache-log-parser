import Head from "next/head"
import { useState } from "react"

import styles from "../styles/Home.module.css"

export default function Home() {
    const [file, setFile] = useState()

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
        <div className={styles.container}>
            <Head>
                <title>Apache Log Parser</title>
                <meta name="Apache Log Parser and Aggregator" />
            </Head>

            <main className={styles.main}>
                <span>Example Apache Common log line</span>
                <p>
                    {
                        '132.128.161.195 - - [25/Jan/2022:20:08:53 -0700] "HEAD /synergize/deploy/cutting-edge/convergence HTTP/2.0" 301 14575'
                    }
                </p>

                <form onSubmit={handleSubmit}>
                    <h1>Apache Log Upload</h1>
                    <input type="file" onChange={handleChange} />
                    <button type="submit">Upload</button>
                </form>
            </main>
        </div>
    )
}
