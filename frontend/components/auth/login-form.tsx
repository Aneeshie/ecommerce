"use client"

import { zodResolver } from "@hookform/resolvers/zod"
import { Controller, useForm } from "react-hook-form"
import { toast } from "sonner"
import { z } from "zod/v4-mini"
import { Card, CardContent, CardFooter, CardHeader, CardTitle, CardDescription } from "../ui/card"
import { Field, FieldError, FieldGroup, FieldLabel } from "../ui/field"
import { Button } from "../ui/button"
import { Input } from "../ui/input"
import { PasswordInput } from "./password-input"
import Link from "next/link"
import { useRouter } from "next/navigation"
import { User } from "@/types/user"
import { useUserStore } from "@/stores/auth-store"

const formSchema = z.object({
  email: z.email(),
  password: z.string()
})

export function LoginForm() {

  const router = useRouter()

  const {setUser} = useUserStore()

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      email: "",
      password: ""
    }
  })

  const onSubmit = async (data: z.infer<typeof formSchema>) => {
    const url = "http://localhost:8080/api/v1/auth/login"

    try {
      const resp = await fetch(url, {
        method: 'POST',
        credentials: "include",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify({email: data.email, password: data.password})
      })

      if (!resp.ok) {
        const err = await resp.json();
        toast.error(err.message ?? "Login failed")
        return;
      }

      const meResp = await fetch("http://localhost:8080/api/v1/auth/me", {
        credentials: "include"
      })

      if (!meResp.ok) {
        const err = await meResp.json();
        toast.error(err.message ?? "failed to get me resp")
        return;
      }

      const user = await meResp.json();

      setUser(user)

      toast.success("Login successfull!")
      router.push("/")
    } catch (err) {
      toast.error("Network error")
    }
  }

    return (
      <Card className="w-full sm:max-w-md shadow-2xl border-white/5 bg-zinc-950/50 backdrop-blur-xl">
        <CardHeader className="space-y-2 pb-6 text-center">
          <CardTitle className="text-2xl font-semibold tracking-tight">Welcome back</CardTitle>
          <CardDescription>Enter your email below to log into your account.</CardDescription>
        </CardHeader>
        <CardContent>
          <form id="form-rhf-login" onSubmit={form.handleSubmit(onSubmit)}>
            <FieldGroup className="space-y-4">
              <Controller
                name="email"
                control={form.control}
                render={({ field, fieldState }) => (
                  <Field data-invalid={fieldState.invalid}>
                    <FieldLabel htmlFor="form-rhf-demo-title">
                      Email
                    </FieldLabel>
                    <Input
                      {...field}
                      id="form-rhf-title"
                      aria-invalid={fieldState.invalid}
                      placeholder="Enter your email address."
                      autoComplete="off"
                    />
                    {fieldState.invalid && (
                      <FieldError errors={[fieldState.error]} />
                    )}
                  </Field>
                )}
            />
            <Controller
              name="password"
              control={form.control}
              render={({ field, fieldState }) => (
                <Field data-invalid={fieldState.invalid}>
                  <FieldLabel htmlFor="form-rhf-demo-title">
                   Password
                  </FieldLabel>
                  <PasswordInput
                    {...field}
                    id="form-rhf-title"
                    aria-invalid={fieldState.invalid}
                    placeholder="Enter your password"
                    autoComplete="off"
                  />
                  {fieldState.invalid && (
                    <FieldError errors={[fieldState.error]} />
                  )}
                </Field>
              )}
            />
            </FieldGroup>
          </form>
        </CardContent>
        <CardFooter className="flex flex-col gap-4 pb-8">
          <Button type="submit" form="form-rhf-login" className="w-full h-11 text-base font-medium shadow-sm transition-all hover:scale-[1.02]">
            Sign In
          </Button>
          <div className="text-center text-sm text-zinc-500 mt-2">
            Don't have an account?{" "}
            <Link href="/signup" className="text-zinc-300 hover:text-white hover:underline transition-colors">
              Sign up
            </Link>
          </div>
        </CardFooter>
      </Card>
    )}
