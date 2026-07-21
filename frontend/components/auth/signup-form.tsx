"use client"

import { zodResolver } from "@hookform/resolvers/zod"
import { Controller, useForm } from "react-hook-form"
import { toast } from "sonner"
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "../ui/card"
import { Field, FieldError, FieldGroup, FieldLabel } from "../ui/field"
import { Button } from "../ui/button"
import { Input } from "../ui/input"
import { PasswordInput } from "./password-input"
import { useRouter } from "next/navigation"
import z from "zod"
import { useUserStore } from "@/stores/auth-store"

const formSchema = z.object({
  name: z.string().min(3, "Name must be at least 3 characters."),
  email: z.email("Invalid email address."),
  password: z
    .string()
    .min(8, "Must be at least 8 characters long")
    .regex(/[A-Z]/, "Must contain at least one uppercase letter")
    .regex(/[a-z]/, "Must contain at least one lowercase letter")
    .regex(/[0-9]/, "Must contain at least one number")
    .regex(/[^A-Za-z0-9]/, "Must contain at least one special character")
})

export function SignUpForm() {

  const router = useRouter()

  const {setUser} = useUserStore()

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    mode: "onChange",
    defaultValues: {
      name: "",
      email: "",
      password: ""
    }
  })

  const onSubmit = async (data: z.infer<typeof formSchema>) => {
    const url = "http://localhost:8080/api/v1/auth/register"
        try {
          const resp = await fetch(url, {
            method: 'POST',
            credentials: "include",
            headers: {
              "Content-Type": "application/json"
            },
            body: JSON.stringify({name: data.name, email: data.email, password: data.password})
          })

          if (!resp.ok) {
            const err = await resp.json();
            toast.error(err.message ?? "SignUp failed")
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

          toast.success("Account created successfully!")
          router.push("/")
        } catch (err) {
          toast.error("Network error")
          console.error(err)
        }
  }

    return (
      <Card className="w-full sm:max-w-md">
        <CardHeader>
          <CardTitle>Login Page</CardTitle>
        </CardHeader>
        <CardContent>
          <form id="form-rhf-signup" onSubmit={form.handleSubmit(onSubmit)}>
            <FieldGroup>
              <Controller
                name="name"
                control={form.control}
                render={({ field, fieldState }) => (
                  <Field data-invalid={fieldState.invalid}>
                    <FieldLabel htmlFor="form-rhf-name">
                     Name
                    </FieldLabel>
                    <Input
                      {...field}
                      id="form-rhf-name"
                      aria-invalid={fieldState.invalid}
                      placeholder="Enter your name."
                      autoComplete="off"
                    />
                    {fieldState.invalid && (
                      <FieldError errors={[fieldState.error]} />
                    )}
                  </Field>
                )}
            />
            <Controller
              name="email"
              control={form.control}
              render={({ field, fieldState }) => (
                <Field data-invalid={fieldState.invalid}>
                  <FieldLabel htmlFor="form-rhf-email">
                    Email
                  </FieldLabel>
                  <Input
                    {...field}
                    id="form-rhf-email"
                    aria-invalid={fieldState.invalid}
                    placeholder="Enter your email"
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
                    <FieldLabel htmlFor="form-rhf-password">
                     Password
                    </FieldLabel>
                    <PasswordInput
                      {...field}
                      id="form-rhf-password"
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
        <CardFooter>
          <Field orientation="horizontal">
            <Button type="button" variant="outline" onClick={() => form.reset()}>
              Reset
            </Button>
            <Button type="submit" form="form-rhf-signup">
              Submit
            </Button>
          </Field>
        </CardFooter>
      </Card>
    )}
