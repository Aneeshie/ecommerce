import { User } from "@/types/user";
import { create } from "zustand";

type UserStore = {
  user: User | null;
  loading: boolean

  setUser: (user: User) => void
  clearUser: () => void
  setLoading: (isLoading: boolean) => void
}

export const useUserStore = create<UserStore>((set) => ({
  user: null,
  loading: true,

  setUser: (user) => set({ user }),
  clearUser: () => set( {user: null}),
  setLoading: (isLoading: boolean) => set({loading: isLoading})
}))
