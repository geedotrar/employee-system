import {
    FormControl,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
  } from "@/components/ui/form"
  import { Input } from "@/components/ui/input"
  import { Control } from "react-hook-form"
  
  interface CustomFormFieldProps {
    name: string
    label: string
    placeholder: string
    control: Control<any>
    type?: string
  }
  
  export function CustomFormField({
    name,
    label,
    placeholder,
    control,
    type = "text",
  }: CustomFormFieldProps) {
    return (
      <FormField
        control={control}
        name={name}
        render={({ field }) => (
          <FormItem>
            <FormLabel>{label}</FormLabel>
            <FormControl>
              <Input type={type} placeholder={placeholder} {...field} />
            </FormControl>
            <FormMessage />
          </FormItem>
        )}
      />
    )
  }