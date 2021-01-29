const speed = 200;
var defaultFetchParam = {
    method: 'POST',
    mode: 'cors',
    cache: 'no-cache',
    headers: {
        "Content-Type": "application/json",
    },
};

var months = [];
var expenses = [];
var saving = [];

var host = "http://localhost:8000/";
//getExpensesOption();

setInterval(() => {
    var today = new Date();
    var date = today.getFullYear() + '-' + (today.getMonth() + 1) + '-' + today.getDate();
    var hours, mins, secs;


    if (today.getHours() < 10) {
        hours = "0" + today.getHours();
    } else {
        hours = today.getHours();
    }

    if (today.getMinutes() < 10) {
        mins = "0" + today.getMinutes();
    } else {
        mins = today.getMinutes();
    }

    if (today.getSeconds() < 10) {
        secs = "0" + today.getSeconds();
    } else {
        secs = today.getSeconds();
    }


    var time = hours + ":" + mins + ":" + secs;
    document.getElementById("date").innerHTML = date;
    document.getElementById("time").innerHTML = time;
}, 1000);



function selectAction(e) {
    e.classList.toggle("unfold");
    var deleteIcon = document.querySelectorAll(".delete-icon");
    var addIcon = document.querySelector(".add-option");

    if (e.classList.contains("save-container") && e.classList.contains("unfold")) {
        document.querySelector(".piggy-bank").style.animation = "piggy 5s linear infinite";
    } else {
        document.querySelector(".piggy-bank").style.animation = "none";
    }


    if (e.id == "editOption") {
        document.getElementById("expensesOptionInput").value = "";
        for (let icon of deleteIcon) {
            icon.classList.toggle("inline-block");
        }
        addIcon.classList.toggle("block");
    }
}

function handleOption(e) {
    var focus = e.id;
    var selectedOptionExpensesTitle = e.getElementsByTagName("label");
    var getSelectedOption = e.getAttribute("data-option");
    var getCurrentSelectOptionCost = e.getAttribute("currentCost");
    var getRemarkText = e.getAttribute("remarks");
    var generateExpensesContainer = document.createElement('div');
    var containerWrap = document.getElementById("containerWrap");

    if (document.getElementById("expensesContainer") != undefined) {
        var expensesContainer = document.getElementById("expensesContainer");
        containerWrap.removeChild(expensesContainer)
    }

    generateExpensesContainer.classList.add("container");
    generateExpensesContainer.classList.add("key-in-value");
    generateExpensesContainer.setAttribute("id", "expensesContainer");

    generateExpensesContainer.innerHTML = `<h3>ENTER YOUR DAILY EXPENSES<br><span class="money-sign">$</span><span
                                                id="optionSelectedExpenses">${selectedOptionExpensesTitle[0].innerHTML}</span><span class="money-sign">$</span></h3>
                                            <div id="keyInContainer">
                                                <label for="keyInAmount">RM</label>
                                                <input type="number" id="keyInExpenses" value="${getCurrentSelectOptionCost}" step="0.01">
                                            </div>
                                            <div style="text-align:left;padding-left:20px; margin-top:15px;" onclick="document.getElementById('remarkDetails').classList.toggle('active')">
                                                <img src="/static/images/remark.png" alt="remark expenses" style="max-width:40px">
                                                <span>Remarks</span>                                    
                                            </div>
                                            <textarea id="remarkDetails" rows="4" cols="20" style="margin:0 auto 10px; display:none;">${getRemarkText}</textarea>
                                            <div style="display: flex; justify-content: space-around;">
                                            <button type="button" class="primary-button" style="padding:5px 20px; margin:0;" id="addExpenses"
                                              optionExpenses="${getSelectedOption}"  onclick="totalExpenses(this)">OK</button>
                                            <button type="button" class="secondary-button" style="padding:5px 10px; margin:0;"
                                                id="cancelExpenses" onclick="totalExpenses(this)">CANCEL</button>
                                            </div>`

    console.log(getCurrentSelectOptionCost);
    console.log(getSelectedOption);
    containerWrap.appendChild(generateExpensesContainer);
}

function totalExpenses(e) {
    var containerWrap = document.getElementById("containerWrap");
    var expensesContainer = document.getElementById("expensesContainer");
    var getFlyMoneyModel = document.querySelector(".fly-animation");

    if (e.id == "cancelExpenses") {
        containerWrap.removeChild(expensesContainer);
    } else if (e.id == "addExpenses") {
        var keyInExpenses = document.getElementById("keyInExpenses");
        var getOptionExpenses = e.getAttribute("optionExpenses");
        var remarkDetails = document.getElementById("remarkDetails");

        if ((keyInExpenses.value == "") || (keyInExpenses.value == undefined)) {
            console.log("empty keyinexpenses");
            getOptionExpensesInput.style.border = "1px solid red";

        } else {
            document.getElementById(getOptionExpenses + "Option").setAttribute("currentCost", keyInExpenses.value);
            document.getElementById(getOptionExpenses + "Option").setAttribute("remarks", remarkDetails.value);
            var getCurrentExpensesAtt = document.querySelectorAll(".option");
            var targetPrice = 0.00;

            for (let expenses of getCurrentExpensesAtt) {
                var currentCost = parseFloat(expenses.getAttribute("currentCost"));
                if (isNaN(currentCost) == false) {
                    targetPrice += currentCost;
                }
            }
            console.log("target price", targetPrice);
            document.getElementById("totalCost").setAttribute("data-target", targetPrice);
            getFlyMoneyModel.style.display = "block";
            getFlyMoneyModel.style.animation = "fly 3s linear infinite";
            updateCount(getFlyMoneyModel);
        }


    }

}

function updateCount(getFlyMoneyModel) {
    var counter = document.getElementById("totalCost")
    const target = +counter.getAttribute('data-target');
    const count = +counter.innerText;
    var getFlyMoneyModel = document.querySelector(".fly-animation");
    // Lower inc to slow and higher to slow
    var inc = target / speed;

    if (inc <= 0.005 && inc > 0) {
        inc = Math.ceil((0.005) * 100) / 100
    }

    // Check if target is reached
    if (count < target) {
        // Add inc to count and output in counter

        counter.innerText = (count + inc).toFixed(2);
        // Call function every ms
        setTimeout(updateCount, 1);
    } else {
        getFlyMoneyModel.style.display = "none";
        getFlyMoneyModel.style.animation = "";
        counter.innerText = target.toFixed(2);
        clearTimeout(updateCount);
        containerWrap.removeChild(expensesContainer);

    }
}


function promptUserActionSubmit(action) {
    var modal = document.getElementById("popUpModal");
    var savingDescription = document.getElementById("savingDescription");
    var resetValueDescription = document.getElementById("resetValueDescription");
    var totalCostDescription = document.getElementById("totalCostDescription");
    var okBtn = document.getElementById("ok");
    if (action == "submitExpenses" || action == "submitSaving" || action == "reset") {
        modal.style.display = "block";
        if (action == "submitSaving") {
            var saveAmount = document.getElementById("saveAmount");
            if (saveAmount.value.trim() == "") {
                modal.style.display = "none";
                saveAmount.style.border = "1px solid red";
                return false;
            } else {
                saveAmount.style.border = "none";
                var getSavingAmount = saveAmount.value;
                savingDescription.style.display = "block";
                totalCostDescription.style.display = "none";
                resetValueDescription.style.display = "none";
                document.getElementById("submitValue").innerHTML = getSavingAmount;
                okBtn.setAttribute("action", "saving");
            }
        } else if (action == "submitExpenses") {
            var totalCost = document.getElementById("totalCost");
            var getTotalCost = totalCost.getAttribute("data-target")
            if (getTotalCost == "" || getTotalCost == "0") {
                modal.style.display = "none";
                totalCost.style.color = "red";
                return false;
            } else {
                totalCost.style.color = "black";
                savingDescription.style.display = "none";
                totalCostDescription.style.display = "block";
                resetValueDescription.style.display = "none";
                document.getElementById("submitExpenses").innerHTML = getTotalCost;
                okBtn.setAttribute("action", "expenses");
            }

        } else if (action == "reset") {
            savingDescription.style.display = "none";
            totalCostDescription.style.display = "none";
            resetValueDescription.style.display = "block";
            okBtn.setAttribute("action", "reset");
        }
    } else if (action == "cancel") {
        modal.style.display = "none";
    }
}

function executeUserAction(e) {
    if (e.getAttribute("action") == "reset") {
        console.log("reset");
        var getAllOption = document.querySelectorAll(".option");
        for (let option of getAllOption) {
            option.setAttribute("currentCost", "0");
        }

        document.getElementById("totalCost").innerHTML = "0.00";
        document.getElementById("popUpModal").style.display = "none";

    } else if (e.getAttribute("action") == "saving") {
        var data = {}
        data.saving = parseFloat(document.getElementById("saveAmount").value);
        defaultFetchParam.body = JSON.stringify(data)

        var url = host + "addSaving";

        fetch(url, defaultFetchParam)
            .then((res) => {
                return res.json();
            })
            .then((result) => {
                console.log("asdasdasd", result.status)
                document.getElementById("saveAmount").value = "";
                console.log("hmmm 222");
                document.getElementById("popUpModal").style.display = "none";
            })

    } else if (e.getAttribute("action") == "expenses") {
        var expensesOption = document.getElementById("expensesOption");
        var getAllExpenses = expensesOption.querySelectorAll('.option');
        console.log("get all expenses >>", getAllExpenses);
        var expensesArr = [];

        for (let expense of getAllExpenses) {
            var jsonExpenses = {};
            if (expense.getAttribute("currentcost") != "" && expense.getAttribute("currentcost") != undefined) {
                jsonExpenses.expensesOption = expense.getAttribute("data-option");
                jsonExpenses.expensesValue = parseFloat(expense.getAttribute("currentcost"));
                jsonExpenses.expensesRemark = expense.getAttribute("remarks");
                expensesArr.push(jsonExpenses);
            }
        }

        console.log("expenses Arr", expensesArr);
        var expenses = {};
        expenses.allExpenses = expensesArr;
        defaultFetchParam.body = JSON.stringify(expenses)
        var url = host + "addExpenses"
        console.table(expenses)
        fetch(url, defaultFetchParam)
            .then((response) => {
                return response.json()
            })
            .then((result) => {
                console.log("result 123", result)

                if (result.status == "00") {
                    document.getElementById("popUpModal").style.display = "none";
                    document.getElementById("totalCost").innerHTML = "";
                    document.getElementById("totalCost").setAttribute("data-target", "0");
                    for (let expense of getAllExpenses) {
                        expenses.setAttribute("currentCost", "0")
                    }
                }
            })

    }

}

function popUpExpensesOptionModal(action, e = "") {
    if (action == "open" || action == "edit") {
        document.getElementById("expensesOptionModal").style.display = "block";
        document.querySelector("body").style.overflow = "hidden";
        var expensesOptionMessage = document.getElementById("expensesOptionMessage");
        var addOptionBtn = document.getElementById("addOptionBtn");
        var delOptionBtn = document.getElementById("delOptionBtn");
        if (action == "open") {
            var expensesOptionInput = document.getElementById("expensesOptionInput");
            expensesOptionMessage.innerHTML = `Add ${expensesOptionInput.value} as expenses option ?`;
            addOptionBtn.style.display = "inline-block";
            delOptionBtn.style.display = "none";

        } else if (action == "edit") {
            var currentOption = e.parentNode;
            console.log("current option 123", currentOption);
            expensesOptionMessage.innerHTML = `delete ${currentOption.innerText} from expenses option ?`;
            delOptionBtn.style.display = "inline-block";
            delOptionBtn.setAttribute("optionId", e.getAttribute("optionId"));
            addOptionBtn.style.display = "none";

        }
    } else if (action == "cancel") {
        document.getElementById("expensesOptionModal").style.display = "none";
        document.querySelector("body").style.overflow = "initial";
    }

}



function addExpensesOption() {

    var expensesOptionInput = document.getElementById("expensesOptionInput");
    expensesJson = {};
    expensesJson.name = expensesOptionInput.value;

    defaultFetchParam.body = JSON.stringify(expensesJson);
    var url = host + "addExpensesOption";

    fetch(url, defaultFetchParam)
        .then((response) => {
            return response.json()
        })
        .then((result) => {
            console.log("result", result);
            document.querySelector("body").style.overflow = "initial";
            document.getElementById("expensesOptionModal").style.display = "none";
            editOption.click();
            getExpensesOption();
        })
}

function getExpensesOption() {
    var url = host + "readExpensesObject";

    fetch(url)
        .then((response) => {
            return response.json()
        })
        .then((result) => {

        })
}

(async () => {
    let fetchExp = await fetch(host + "readExpensesObject")
    console.log("fetchExp", fetchExp)
    let result = await fetchExp.json()
    console.log("fetchResponse", result)
    var expensesOption = document.getElementById("expensesOption");
    var totalCost = document.getElementById("totalCost");
    totalCost.innerHTML = "";
    expensesOption.innerHTML = "";
    var totalValue = 0.0;
    console.log("expenses option data", result)
    for (var i = 0; i < result.length; i++) {
        var expensesContainer = document.createElement("div");
        expensesContainer.id = result[i].name.toLowerCase() + "Option"
        expensesContainer.setAttribute("data-option", result[i].name.toLowerCase());
        expensesContainer.setAttribute("currentcost", result[i].currentValue);
        expensesContainer.setAttribute("remarks", result[i].currentRemark);
        expensesContainer.setAttribute("onclick", "handleOption(this)")
        expensesContainer.classList.add("option");

        expensesContainer.innerHTML = `<img src="/static/images/delete.svg" width=30 class="delete-icon" onclick="popUpExpensesOptionModal('edit',this)" optionId="${result[i].id}">
                                       <input type="radio" id="${result[i].name.toLowerCase()}" name="expense" value="${result[i].name.toLowerCase()}">
                                       <label for="${result[i].name.toLowerCase()}">${capitalizeWord(result[i].name)}</label>
                                       <div class="selected"></div>
                                       `

        console.log(result[i].currentValue, isNaN(result[i].currentValue), typeof (result[i].currentValue));
        if ((isNaN(result[i].currentValue) == false) && (result[i].currentValue != "0") && (result[i].currentValue != "")) {
            totalValue += parseFloat(result[i].currentValue);
        }
        expensesOption.appendChild(expensesContainer);
        console.log("total Value", totalValue)
    }
    totalCost.innerText = totalValue;
})()

function deleteExpensesOption(e) {

    var optionId = e.getAttribute("optionId");
    console.log("option id", optionId);
    var url = host + "deleteExpensesOption/" + optionId;
    defaultFetchParam.body = JSON.stringify("");

    fetch(url, defaultFetchParam)
        .then((response) => {
            return response.json();
        })
        .then((result) => {
            console.log("result succesfull ?", result);
            document.getElementById("expensesOptionModal").style.display = "none";
            getExpensesOption();
        })

}

function capitalizeWord(word) {
    return word.charAt(0).toUpperCase() + word.slice(1).toLowerCase();
}

function closeExpensesOptionModal() {
    document.getElementById("expensesOptionModal").style.display = "none";
}

function closePopUpModal() {
    document.getElementById("popUpModal").style.display = "none";
}


function getExpensesReport() {
    var url = host + "getFinancialStatistic";
    fetch(url)
        .then((response) => {
            return response.json();
        })
        .then((result) => {
            console.log("statistic data", result);
            for (var i = 0; i < result.report.length; i++) {

                months.push(result.report[i].month);
                expenses.push(result.report[i].totalExpenses);
                saving.push(result.report[i].saving);
            }

        })
        .then(() => {
            var ctx = document.getElementById('myChart').getContext('2d');
            var mixedChart = new Chart(ctx, {
                type: 'line',
                data: {
                    datasets: [{
                        label: 'Saving',
                        data: saving,
                        fill: false,
                        borderColor: "#4DFF33",
                        backgroundColor: "#2EFF0F"

                    }, {
                        label: 'Expenses',
                        data: expenses,
                        fill: false,
                        borderColor: "#FF2C0F",
                        backgroundColor: "#FF1F00"
                    }],
                    labels: months
                },
                options: {
                    scales: {
                        xAxes: [{
                            display: true,
                            scaleLabel: {
                                display: true,
                                labelString: 'Month'
                            }
                        }],
                        yAxes: [{
                            display: true,
                            scaleLabel: {
                                display: true,
                                labelString: '$$'
                            }
                        }]
                    },
                    title: {
                        display: true,
                        text: 'Monthly Money$$ Statistic',
                        fontSize: 16
                    }
                }
            });

            document.getElementById("myChart").onclick = function (evt) {
                var activePoints = mixedChart.getElementsAtEvent(evt);
                console.log("activePoints", activePoints);
                if (activePoints[0]) {
                    var chartData = activePoints[0]['_chart'].config.data;
                    console.log("chartData", chartData);

                    var idx = activePoints[0]['_index'];
                    var dIndex = mixedChart.getDatasetAtEvent(evt)[0]._datasetIndex;

                    var label = chartData.labels[idx];
                    var value = chartData.datasets[dIndex].data[idx];
                    generateExpensesSummary(label)
                    console.log("label", label, "value", value);


                }
            }
        })
}


function generateExpensesSummary(month) {
    var url = host + "generateExpensesSummary/" + month;
    var expensesObject = [];
    var expensesValue = [];
    fetch(url)
        .then((response) => {
            return response.json()
        })
        .then((result) => {
            console.log("result", result);
            for (let option of result) {
                console.log("option");
                expensesObject.push(option.expenses);
                expensesValue.push(option.value)

            }
            console.log("object", expensesObject);
            console.log("value", expensesValue);

        })
        .then(() => {
            var expeneseSummary = document.getElementById('expeneseSummary');
            expeneseSummary.innerHTML = "";
            var summaryCanvas = document.createElement('canvas');
            var ctx = summaryCanvas.getContext('2d');
            var myChart = new Chart(ctx, {
                type: 'pie',
                data: {
                    labels: expensesObject,
                    datasets: [{
                        backgroundColor: [
                            "#2ecc71",
                            "#3498db",
                            "#95a5a6",
                            "#9b59b6",
                            "#f1c40f",
                        ],
                        data: expensesValue
                    }]
                }
            });

            expeneseSummary.appendChild(summaryCanvas);

        })
}
getExpensesReport();
// var options = {
//     threshold: 1,
//     rootMargin: "0px 0px -100px 0px"
// }

// const chartObserver = new IntersectionObserver((entries) => {
//     entries.forEach((entry) => {
//         if (entry.isIntersecting) {

//         }
//     })
// }, options);

// chartObserver.observe(document.getElementById("addSectionId"));