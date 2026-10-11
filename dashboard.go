
package main

import (
	"fmt"
	"net/http"
)

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	fmt.Fprint(w, `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>FinOps Control Center</title>

    <style>
        * {
            box-sizing: border-box;
        }

        body {
            margin: 0;
            padding: 30px 16px;
            font-family: Arial, sans-serif;
            background: #f1f5f9;
            color: #172033;
        }

        .container {
            max-width: 780px;
            margin: auto;
        }

        h1 {
            margin-bottom: 8px;
            color: #2457c5;
        }

        .subtitle {
            color: #64748b;
            line-height: 1.5;
        }

        .card {
            padding: 24px;
            margin-top: 22px;
            background: white;
            border: 1px solid #e2e8f0;
            border-radius: 14px;
            box-shadow: 0 4px 15px #00000008;
        }

        .grid {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 16px;
        }

        label {
            display: block;
            margin: 14px 0 7px;
            font-weight: bold;
            font-size: 14px;
        }

        input {
            width: 100%;
            padding: 12px;
            border: 1px solid #cbd5e1;
            border-radius: 8px;
            font-size: 16px;
        }

        button {
            width: 100%;
            margin-top: 22px;
            padding: 14px;
            background: #2457c5;
            color: white;
            border: none;
            border-radius: 8px;
            font-size: 16px;
            font-weight: bold;
            cursor: pointer;
        }

        button:hover {
            background: #1946a7;
        }

        button:disabled {
            opacity: 0.6;
            cursor: wait;
        }

        .summary {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 16px;
        }

        .metric {
            padding: 18px;
            background: #f8fafc;
            border-radius: 10px;
        }

        .metric-label {
            color: #64748b;
            font-size: 13px;
        }

        .metric-value {
            margin-top: 8px;
            font-size: 25px;
            font-weight: bold;
            overflow-wrap: anywhere;
        }

        #result {
            display: none;
            margin-top: 20px;
            padding: 18px;
            border-radius: 10px;
            line-height: 1.6;
            overflow-wrap: anywhere;
        }

        .allowed {
            display: block !important;
            background: #dcfce7;
            color: #166534;
        }

        .rejected {
            display: block !important;
            background: #fee2e2;
            color: #991b1b;
        }

        .error {
            display: block !important;
            background: #fef3c7;
            color: #92400e;
        }

        .muted {
            color: #64748b;
            font-size: 13px;
        }

        @media (max-width: 520px) {
            .grid,
            .summary {
                grid-template-columns: 1fr;
            }

            .card {
                padding: 18px;
            }
        }
    </style>
</head>

<body>
    <main class="container">
        <h1>FinOps Control Center</h1>
        <p class="subtitle">
            Estimate Kubernetes workload costs and check whether they
            fit within your monthly budget.
        </p>

        <section class="card">
            <h2>Workload calculator</h2>
            <p class="muted">
                Example pricing: CPU ₹2 per core-hour and memory
                ₹0.50 per GiB-hour, using 730 hours per month.
            </p>

            <form id="calculator">
                <div class="grid">
                    <div>
                        <label for="cpu">CPU cores</label>
                        <input id="cpu" type="number"
                            min="0.01" step="any" value="2" required>
                    </div>

                    <div>
                        <label for="memory">Memory (GiB)</label>
                        <input id="memory" type="number"
                            min="0.01" step="any" value="4" required>
                    </div>

                    <div>
                        <label for="replicas">Number of replicas</label>
                        <input id="replicas" type="number"
                            min="1" step="1" value="3" required>
                    </div>

                    <div>
                        <label for="budget">Monthly budget (₹)</label>
                        <input id="budget" type="number"
                            min="0" step="any" value="10000" required>
                    </div>
                </div>

                <button id="calculateButton" type="submit">
                    Calculate and Check Budget
                </button>
            </form>

            <div id="result" role="status" aria-live="polite"></div>
        </section>

        <section class="card">
            <h2>How the system works</h2>
            <p>1. Enter the workload's CPU and memory requirements.</p>
            <p>2. Set the number of replicas and the monthly budget.</p>
            <p>3. The Go backend calculates the estimated monthly cost.</p>
            <p>4. The dashboard shows whether the budget check passed.</p>
            <p class="muted">
                This calculator checks your budget through the Go API.
                It does not create a Kubernetes Deployment.
            </p>
        </section>
    </main>

    <script>
        const form = document.getElementById('calculator');
        const resultBox = document.getElementById('result');
        const calculateButton = document.getElementById('calculateButton');

        form.addEventListener('submit', async function (event) {
            event.preventDefault();

            const cpu = Number(document.getElementById('cpu').value);
            const memory = Number(document.getElementById('memory').value);
            const replicas = Number(document.getElementById('replicas').value);
            const budget = Number(document.getElementById('budget').value);

            resultBox.className = 'error';
            resultBox.textContent = '';
            resultBox.style.display = 'block';

            if (!Number.isFinite(cpu) || cpu <= 0 ||
                !Number.isFinite(memory) || memory <= 0 ||
                !Number.isInteger(replicas) || replicas < 1 ||
                !Number.isFinite(budget) || budget < 0) {
                resultBox.textContent =
                    'Enter valid CPU, memory, replicas and budget values.';
                return;
            }

            calculateButton.disabled = true;
            calculateButton.textContent = 'Calculating...';
            resultBox.textContent = 'Contacting the Go backend...';

            try {
                const response = await fetch('/evaluate', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        CPUCores: cpu,
                        MemoryGiB: memory,
                        Replicas: replicas,
                        Budget: budget
                    })
                });

                if (!response.ok) {
                    throw new Error('Server returned HTTP ' + response.status);
                }

                const data = await response.json();
                const cost = Number(data.EstimatedMonthlyCost);

                resultBox.className = data.Allowed ? 'allowed' : 'rejected';

                const heading = document.createElement('h3');
                heading.textContent = data.Allowed
                    ? 'ALLOW — Within Budget'
                    : 'REJECT — Budget Check Failed';

                const costLine = document.createElement('p');
                costLine.textContent = 'Estimated monthly cost: ₹' +
                    (Number.isFinite(cost) ? cost.toFixed(2) : '0.00');

                const budgetLine = document.createElement('p');
                budgetLine.textContent =
                    'Monthly budget: ₹' + budget.toFixed(2);

                const reasonLine = document.createElement('p');
                reasonLine.textContent = 'Reason: ' + (data.Reason || 'No reason provided.');

                resultBox.replaceChildren(
                    heading,
                    costLine,
                    budgetLine,
                    reasonLine
                );
            } catch (error) {
                resultBox.className = 'error';
                resultBox.textContent =
                    'Could not contact the Go backend. ' + error.message;
            } finally {
                calculateButton.disabled = false;
                calculateButton.textContent = 'Calculate and Check Budget';
            }
        });
    </script>
</body>
</html>
`)
}
