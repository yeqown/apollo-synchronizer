/* 
    @param {number} num
    @return {string}
    @example
    formatNumber(123456789) // '123,456,789'
    formatNumber(0) // '-'
    */
const formatNumber = (num) => {
    if (typeof num !== 'number' || isNaN(num)) {
        return "NaN";
    }

    if (num <= 0) {
        return "--";
    }

    return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
}

export { formatNumber }