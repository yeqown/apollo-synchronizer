/*
    * @param {number} ts
    * @return {string}
    * @example
    * formatTs(1547782400) // "2019-01-01 00:00:00"
    */
const formatTs = (ts) => {
    const date = new Date(ts);

    // parse ts in format "2019-01-01T00:00:00.000Z"
    return `${date.getFullYear()}-${date.getMonth() + 1}-${date.getDate()} ${date.getHours()}:${date.getMinutes()}:${date.getSeconds()}`;
}

/*
    * @param {number} seconds
    * @return {string}
    * @example
    * humanizeTime(1547782400) // "1h 0m 0s"
    */
const humanizeTime = (seconds) => {
    const secs = Math.floor(seconds % 60);
    const mins = Math.floor((seconds / 60) % 60);
    const hrs = Math.floor(seconds / 3600);

    return `${hrs}h ${mins}m ${secs}s`;
}

export { humanizeTime, formatTs };

