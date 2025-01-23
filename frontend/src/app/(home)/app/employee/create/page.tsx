"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { CustomFormField } from "@/components/ui/custom-form";

const formSchema = z.object({
  username: z.string().min(2, {
    message: "Username must be at least 2 characters.",
  }),
  fullName: z.string().min(3, {
    message: "Fullname must be at least 2 characters.",
  }),
  email: z.string().email({
    message: "Email not valid.",
  }),
  phoneNumber: z.string().min(10, {
    message: "Phone must be at least 10 characters.",
  }),
  position: z.string().min(2, {
    message: "Position name must be at least 2 characters",
  }),
});

export function CreateForm() {
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      username: "",
      fullName: "",
      email: "",
      phoneNumber: "",
      position: "",
    },
  });

  const formFields = [
    { name: "username", label: "Username", placeholder: "alfiankafilah" },
    {
      name: "fullName",
      label: "Full Name",
      placeholder: "Alfian Kafilah Baits",
    },
    {
      name: "email",
      label: "Email",
      placeholder: "alfian@example.com",
      type: "email",
    },
    { name: "phoneNumber", label: "Phone Number", placeholder: "08123456789" },
    { name: "position", label: "Position", placeholder: "CTO" },
  ];

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(console.log)} className="space-y-4">
        {formFields.map((field) => (
          <CustomFormField key={field.name} control={form.control} {...field} />
        ))}
        <Button type="submit">Submit</Button>
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
