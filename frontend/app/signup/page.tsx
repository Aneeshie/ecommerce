import { LoginForm } from "@/components/auth/login-form";
import { SignUpForm } from "@/components/auth/signup-form";

export default function SignupPage() {

  return (
      <main className="relative flex min-h-screen items-center justify-center overflow-hidden bg-black px-4">
        <div className="absolute inset-0 bg-[radial-gradient(ellipse_at_top,_var(--tw-gradient-stops))] from-zinc-800/20 via-black to-black -z-10" />
        <SignUpForm />
      </main>
  )
}
