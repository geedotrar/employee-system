import { Employee, columns } from "./columns"
import { DataTable } from "./data-table"

async function getData(): Promise<Employee[]> {
  // Fetch data from your API here.
  return [
    {
      id: "728ed52f",
      name: "Alfian",
      employee_code: 1234,
      position: "CTO",
    },
    // ...
  ]
}

export default async function EmployeePage() {
  const data = await getData()

  return (
    <div className="container mx-auto py-10">
      <DataTable columns={columns} data={data} />
    </div>
  )
}
