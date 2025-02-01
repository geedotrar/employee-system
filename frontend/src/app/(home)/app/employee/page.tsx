import { Button } from "@/components/ui/button";
import { Employee, columns } from "./columns";
import { DataTable } from "./data-table";
import Link from "next/link";
import { CirclePlusIcon } from "lucide-react";

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
  ];
}

export default async function EmployeePage() {
  const data = await getData();

  return (
    <div className="container mx-auto">
      <div className="mb-3 flex justify-end">
        <Button asChild>
          <Link href="/app/employee/create">
            {" "}
            <CirclePlusIcon /> Add Data
          </Link>
        </Button>
      </div>
      <DataTable columns={columns} data={data} />{" "}
    </div>
  );
}
