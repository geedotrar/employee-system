import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

export default function Login() {
  return (
    <div className="flex items-center justify-center h-screen">
      <div className="w-full max-w-xs space-y-4">
        <h1 className="text-center text-2xl">Login Page</h1>
        <div>
          <label className="block text-left">Employee Code</label>
          <Input name="employee_code" type="text" placeholder="Employee Code" />
        </div>
        <div>
          <label className="block text-left">Password</label>
          <Input name="password" type="password" placeholder="Password" />
        </div>
        <div className="flex justify-center">
          <Button className="mt-1 w-full">Login</Button>
        </div>
      </div>
    </div>
  );
}
