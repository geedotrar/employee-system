"use client"

import { ColumnDef } from "@tanstack/react-table"

// This type is used to define the shape of our data.
// You can use a Zod schema here if you want.
export type Employee = {
  id: string
  name: string
  employee_code: number
  position: string
}

export const columns: ColumnDef<Employee>[] = [
  {
    accessorKey: "name",
    header: "Name",
  },
  {
    accessorKey: "employee_code",
    header: "Employee Code",
  },
  {
    accessorKey: "position",
    header: "Position",
  },
]
