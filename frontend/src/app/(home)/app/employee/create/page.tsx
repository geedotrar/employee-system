"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { EmployeeSchema } from "@/schemas";
import { Input } from "@/components/ui/input";
import Link from "next/link";

export function CreateForm() {
  const form = useForm<z.infer<typeof EmployeeSchema>>({
    resolver: zodResolver(EmployeeSchema),
    defaultValues: {
      username: "",
      fullName: "",
      email: "",
      phoneNumber: "",
      position: "",
    },
  });

  const formFields = [
    {
      name: "username",
      label: "Username",
      placeholder: "alfiankafilah",
      type: "text",
    },
    {
      name: "fullName",
      label: "Full Name",
      placeholder: "Alfian Kafilah Baits",
      type: "text",
    },
    {
      name: "email",
      label: "Email",
      placeholder: "alfian@example.com",
      type: "email",
    },
    {
      name: "phoneNumber",
      label: "Phone Number",
      placeholder: "08123456789",
      type: "text",
    },
    { name: "position", label: "Position", placeholder: "CTO", type: "text" },
  ] as const;

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(console.log)} className="space-y-4">
        {formFields.map((item) => (
          <FormField
            key={item.name}
            control={form.control}
            name={item.name}
            render={({ field }) => (
              <FormItem>
                <FormLabel>{item.label}</FormLabel>
                <FormControl>
                  <Input
                    type={item.type}
                    placeholder={item.placeholder}
                    {...field}
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
        ))}
        <div className="flex items-center justify-end gap-x-2">
          <Button variant="outline" type="button" asChild>
            <Link href="/app/employee">Back</Link>
          </Button>
          <Button type="submit">Submit</Button>
        </div>
      </form>
    </Form>
  );
}

export default function CreateEmployeePage() {
  return (
    <div className="container mx-auto">
      <CreateForm />
    </div>
  );
}
