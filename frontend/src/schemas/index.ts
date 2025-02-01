import { z } from "zod";

export const EmployeeSchema = z.object({
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
