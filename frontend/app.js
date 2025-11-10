const baseUrl = 'http://localhost:8080/api/v1'

let cachedEmployees = []

async function fetchEmployees() {
  const listMessage = document.getElementById('listMessage')
  if (listMessage) listMessage.textContent = 'Loading...'
  try {
    const res = await fetch(`${baseUrl}/employees`)
    if (!res.ok) throw new Error(`Server returned ${res.status}`)
    const data = await res.json()
    cachedEmployees = data || []
    renderEmployees(cachedEmployees)
    populateEmployeeSelects(cachedEmployees)
    if (listMessage) listMessage.textContent = ''
  } catch (err) {
    if (listMessage) listMessage.textContent = 'Failed to load employees: ' + err.message
    renderEmployees([])
  }
}

function renderEmployees(employees) {
  const tbody = document.querySelector('#employeesTable tbody')
  if (!tbody) return
  tbody.innerHTML = ''
  if (!employees || employees.length === 0) {
    const tr = document.createElement('tr')
    const td = document.createElement('td')
    td.colSpan = 7 // Updated colspan
    td.textContent = 'No employees found.'
    tr.appendChild(td)
    tbody.appendChild(tr)
    return
  }

  employees.forEach(emp => {
    const tr = document.createElement('tr')
    tr.innerHTML = `
      <td>${emp.id ?? ''}</td>
      <td>${escapeHtml(emp.name)}</td>
      <td>${escapeHtml(emp.position)}</td>
      <td>${emp.base_salary}</td>
      <td>${emp.allowance}</td>
      <td>${new Date(emp.created_at).toLocaleString()}</td>
      <td><button class="editBtn" data-id="${emp.id}">Edit</button></td>
    `
    tbody.appendChild(tr)
  })
}

function escapeHtml(s) {
  if (!s) return ''
  return s.replace(/[&<>"']/g, (c) => ({
    '&': '&amp;',
    '<': '&lt;',
    '>': '&gt;',
    '"': '&quot;',
    "'": '&#39;'
  })[c])
}

async function createEmployee(payload) {
  const fm = document.getElementById('formMessage')
  if (fm) fm.textContent = 'Creating...'
  try {
    const res = await fetch(`${baseUrl}/employees`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    })
    if (res.status === 201) {
      const created = await res.json()
      if (fm) fm.textContent = 'Employee created (ID: ' + (created.id ?? '?') + ')'
      document.getElementById('employeeForm').reset()
      await fetchEmployees()
    } else {
      const errBody = await res.json().catch(() => ({}))
      if (fm) fm.textContent = 'Failed to create: ' + (errBody.error || res.status)
    }
  } catch (err) {
    if (fm) fm.textContent = 'Failed to create: ' + err.message
  }
}

async function updateEmployee(id, payload) {
  const fm = document.getElementById('formMessage')
  if (fm) fm.textContent = 'Updating...'
  try {
    const res = await fetch(`${baseUrl}/employees/${id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    })
    if (res.ok) {
      const updated = await res.json()
      if (fm) fm.textContent = 'Employee updated (ID: ' + (updated.id ?? '?') + ')'
      document.getElementById('employeeForm').reset()
      document.getElementById('employee_id').value = ''
      document.querySelector('#employeeForm button[type="submit"]').textContent = 'Create'; // Reset button text
      await fetchEmployees()
    } else {
      const errBody = await res.json().catch(() => ({}))
      if (fm) fm.textContent = 'Failed to update: ' + (errBody.error || res.status)
    }
  } catch (err) {
    if (fm) fm.textContent = 'Failed to update: ' + err.message
  }
}

function readForm() {
  const name = document.getElementById('name').value.trim()
  const base_salary = parseFloat(document.getElementById('base_salary').value)
  const allowance = parseFloat(document.getElementById('allowance').value)
  const position = document.getElementById('position').value.trim()
  return { name, base_salary, allowance, position }
}

function validate(payload) {
  if (!payload.name) return 'Name is required'
  if (!payload.position) return 'Position is required'
  if (Number.isNaN(payload.base_salary)) return 'Base salary must be a number'
  if (Number.isNaN(payload.allowance)) return 'Allowance must be a number'
  return null
}

function populateEmployeeSelects(employees) {
  const selects = [
    document.getElementById('attendance_employee'),
    document.getElementById('payroll_employee'),
    document.getElementById('attendance_history_employee')
  ]
  
  selects.forEach(sel => {
    if (!sel) return
    // clear existing options except placeholder
    sel.querySelectorAll('option:not([value=""])').forEach(n => n.remove())
    employees.forEach(e => {
      const opt = document.createElement('option')
      opt.value = e.id
      opt.textContent = `${e.id} — ${e.name}`
      sel.appendChild(opt)
    })
  })
}

async function recordAttendance(payload) {
  const fm = document.getElementById('attendanceMessage')
  if (fm) fm.textContent = 'Recording...'
  try {
    const res = await fetch(`${baseUrl}/attendances`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    })
    if (res.status === 201) {
      const created = await res.json()
      if (fm) fm.textContent = 'Attendance recorded (ID: ' + (created.id ?? '?') + ')'
      document.getElementById('attendanceForm').reset()
    } else {
      const errBody = await res.json().catch(() => ({}))
      if (fm) fm.textContent = 'Failed to record: ' + (errBody.error || res.status)
    }
  } catch (err) {
    if (fm) fm.textContent = 'Failed to record: ' + err.message
  }
}

async function recordCheckout(payload) {
  const fm = document.getElementById('attendanceMessage')
  if (fm) fm.textContent = 'Recording checkout...'
  try {
    const res = await fetch(`${baseUrl}/attendances/checkout`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    })
    if (res.ok) { // status 200
      const updated = await res.json()
      if (fm) fm.textContent = 'Checkout recorded (ID: ' + (updated.id ?? '?') + ')'
      document.getElementById('attendanceForm').reset()
    } else {
      const errBody = await res.json().catch(() => ({}))
      if (fm) fm.textContent = 'Failed to record checkout: ' + (errBody.error || res.status)
    }
  } catch (err) {
    if (fm) fm.textContent = 'Failed to record checkout: ' + err.message
  }
}

async function fetchAttendanceByPeriod(employeeID, from, to) {
  const listMessage = document.getElementById('attendanceHistoryMessage')
  if (listMessage) listMessage.textContent = 'Loading...'
  try {
    const res = await fetch(`${baseUrl}/attendances?employee_id=${employeeID}&from=${from}&to=${to}`)
    if (!res.ok) throw new Error(`Server returned ${res.status}`)
    const data = await res.json()
    renderAttendanceHistory(data)
    if (listMessage) listMessage.textContent = ''
  } catch (err) {
    if (listMessage) listMessage.textContent = 'Failed to load attendance history: ' + err.message
    renderAttendanceHistory([])
  }
}

function renderAttendanceHistory(attendances) {
  const tbody = document.querySelector('#attendanceHistoryTable tbody')
  tbody.innerHTML = ''
  if (!attendances || attendances.length === 0) {
    const tr = document.createElement('tr')
    const td = document.createElement('td')
    td.colSpan = 5
    td.textContent = 'No attendance records found for this period.'
    tr.appendChild(td)
    tbody.appendChild(tr)
    return
  }
  attendances.forEach(att => {
    const tr = document.createElement('tr')
    tr.innerHTML = `
      <td>${att.id ?? ''}</td>
      <td>${new Date(att.date).toLocaleDateString()}</td>
      <td>${escapeHtml(att.status)}</td>
      <td>${att.check_in ? new Date(att.check_in).toLocaleTimeString() : 'N/A'}</td>
      <td>${att.check_out ? new Date(att.check_out).toLocaleTimeString() : 'N/A'}</td>
    `
    tbody.appendChild(tr)
  })
}

async function generatePayroll(payload) {
  const fm = document.getElementById('payrollMessage')
  if (fm) fm.textContent = 'Generating...'
  try {
    const res = await fetch(`${baseUrl}/payroll/generate`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    })
    if (res.status === 201) {
      const created = await res.json()
      if (fm) fm.textContent = 'Payroll generated (ID: ' + (created.id ?? '?') + ')'
      document.getElementById('payrollForm').reset()
      await fetchPayrollSlips()
    } else {
      const errBody = await res.json().catch(() => ({}))
      if (fm) fm.textContent = 'Failed to generate: ' + (errBody.error || res.status)
    }
  } catch (err) {
    if (fm) fm.textContent = 'Failed to generate: ' + err.message
  }
}

async function fetchPayrollSlips() {
  const listMessage = document.getElementById('payrollListMessage')
  if (listMessage) listMessage.textContent = 'Loading...'
  try {
    const res = await fetch(`${baseUrl}/payroll/slips`)
    if (!res.ok) throw new Error(`Server returned ${res.status}`)
    const data = await res.json()
    renderPayrollSlips(data)
    if (listMessage) listMessage.textContent = ''
  } catch (err) {
    if (listMessage) listMessage.textContent = 'Failed to load payroll slips: ' + err.message
    renderPayrollSlips([])
  }
}

function renderPayrollSlips(slips) {
  const tbody = document.querySelector('#payrollTable tbody')
  tbody.innerHTML = ''
  if (!slips || slips.length === 0) {
    const tr = document.createElement('tr')
    const td = document.createElement('td')
    td.colSpan = 5
    td.textContent = 'No payroll slips found.'
    tr.appendChild(td)
    tbody.appendChild(tr)
    return
  }
  slips.forEach(s => {
    const tr = document.createElement('tr')
    tr.innerHTML = `
      <td>${s.id ?? ''}</td>
      <td>${s.employee_id ?? ''}</td>
      <td>${new Date(s.period).toLocaleDateString()}</td>
      <td>${s.take_home_pay}</td>
      <td>${new Date(s.generated_at).toLocaleString()}</td>
    `
    tbody.appendChild(tr)
  })
}

function switchView(view) {
  document.querySelectorAll('.view').forEach(el => el.style.display = 'none')
  document.querySelectorAll('.tabs button').forEach(b => b.classList.remove('active'))
  document.getElementById('view-' + view).style.display = ''
  document.querySelector(`.tabs button[data-view="${view}"]`).classList.add('active')
}

document.addEventListener('DOMContentLoaded', () => {
  // Set date fields to today
  const today = new Date();
  const yyyy = today.getFullYear();
  const mm = String(today.getMonth() + 1).padStart(2, '0'); // Months are 0-based
  const dd = String(today.getDate()).padStart(2, '0');
  const todayString = `${yyyy}-${mm}-${dd}`;
  
  document.getElementById('attendance_from').value = todayString;
  document.getElementById('attendance_to').value = todayString;

  fetchEmployees()
  fetchPayrollSlips()

  // Navigation
  document.querySelectorAll('.tabs button').forEach(btn => {
    btn.addEventListener('click', (e) => {
      const view = e.currentTarget.getAttribute('data-view')
      switchView(view)
    })
  })

  // Employee form
  document.getElementById('employeeForm').addEventListener('submit', async (e) => {
    e.preventDefault()
    const payload = readForm()
    const v = validate(payload)
    const fm = document.getElementById('formMessage')
    if (v) {
      fm.textContent = v
      return
    }
    
    const employeeId = document.getElementById('employee_id').value;
    if (employeeId) {
      await updateEmployee(employeeId, payload);
    } else {
      await createEmployee(payload);
    }
  })
  document.getElementById('refreshBtn').addEventListener('click', () => fetchEmployees())

  // Event listener for edit buttons (using event delegation)
  document.getElementById('employeesTable').addEventListener('click', (e) => {
    if (e.target.classList.contains('editBtn')) {
      const id = e.target.getAttribute('data-id');
      const employee = cachedEmployees.find(emp => emp.id == id);
      if (employee) {
        document.getElementById('employee_id').value = employee.id;
        document.getElementById('name').value = employee.name;
        document.getElementById('base_salary').value = employee.base_salary;
        document.getElementById('allowance').value = employee.allowance;
        document.getElementById('position').value = employee.position;
        
        // Change button text to "Update"
        document.querySelector('#employeeForm button[type="submit"]').textContent = 'Update';
        
        // Scroll to form
        document.getElementById('employeeForm').scrollIntoView({ behavior: 'smooth' });
      }
    }
  });

  // Attendance form
  document.getElementById('checkinBtn').addEventListener('click', async () => {
    const empId = parseInt(document.getElementById('attendance_employee').value)
    if (!empId) {
      document.getElementById('attendanceMessage').textContent = 'Employee is required.'
      return
    }
    const today = new Date();
    const date = `${today.getFullYear()}-${String(today.getMonth() + 1).padStart(2, '0')}-${String(today.getDate()).padStart(2, '0')}T00:00:00Z`;
    const payload = {
      employee_id: empId,
      date: date,
      status: 'PRESENT',
      check_in: today.toISOString()
    }
    await recordAttendance(payload)
  })

  document.getElementById('checkoutBtn').addEventListener('click', async () => {
    const empId = parseInt(document.getElementById('attendance_employee').value)
    if (!empId) {
      document.getElementById('attendanceMessage').textContent = 'Employee is required.'
      return
    }
    const payload = {
      employee_id: empId,
    }
    await recordCheckout(payload)
  })

  document.getElementById('absentBtn').addEventListener('click', async () => {
    const empId = parseInt(document.getElementById('attendance_employee').value)
    if (!empId) {
      document.getElementById('attendanceMessage').textContent = 'Employee is required.'
      return
    }
    const today = new Date();
    const date = `${today.getFullYear()}-${String(today.getMonth() + 1).padStart(2, '0')}-${String(today.getDate()).padStart(2, '0')}T00:00:00Z`;
    const payload = { employee_id: empId, date: date, status: 'ABSENT' }
    await recordAttendance(payload)
  })

  document.getElementById('leaveBtn').addEventListener('click', async () => {
    const empId = parseInt(document.getElementById('attendance_employee').value)
    if (!empId) {
      document.getElementById('attendanceMessage').textContent = 'Employee is required.'
      return
    }
    const today = new Date();
    const date = `${today.getFullYear()}-${String(today.getMonth() + 1).padStart(2, '0')}-${String(today.getDate()).padStart(2, '0')}T00:00:00Z`;
    const payload = { employee_id: empId, date: date, status: 'LEAVE' }
    await recordAttendance(payload)
  })

  // Attendance History form
  document.getElementById('attendanceHistoryForm').addEventListener('submit', async (e) => {
    e.preventDefault()
    const empId = document.getElementById('attendance_history_employee').value
    const from = document.getElementById('attendance_from').value
    const to = document.getElementById('attendance_to').value
    if (!empId || !from || !to) {
      document.getElementById('attendanceHistoryMessage').textContent = 'Employee, from date, and to date are required.'
      return
    }
    await fetchAttendanceByPeriod(empId, from, to)
  })

  // Payroll form
  document.getElementById('payrollForm').addEventListener('submit', async (e) => {
    e.preventDefault()
    const emp = document.getElementById('payroll_employee').value
    const periodMonth = document.getElementById('payroll_period').value
    // convert month input (YYYY-MM) to YYYY-MM-01
    const period = periodMonth ? periodMonth + '-01' : ''
    const payload = { employee_id: emp ? parseInt(emp) : 0, period: period }
    await generatePayroll(payload)
  })
  document.getElementById('refreshPayroll').addEventListener('click', () => fetchPayrollSlips())
})
