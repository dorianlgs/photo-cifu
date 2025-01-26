import { fail, error, redirect } from "@sveltejs/kit"

/** @type {import('./$types').Actions} */
export const actions = {
    signIn: async ({ request, locals }) => {
        const formData = await request.formData()
        const errors: { [fieldName: string]: string } = {}

        const email = formData.get("email")?.toString() ?? ""
        if (email.length < 6) {
            errors["email"] = "Email is required"
        } else if (email.length > 500) {
            errors["email"] = "Email too long"
        } else if (!email.includes("@") || !email.includes(".")) {
            errors["email"] = "Invalid email"
        }

        const password = formData.get("password")?.toString() ?? ""
        if (password.length > 500) {
            errors["password"] = "Password too long"
        }

        if (Object.keys(errors).length > 0) {
            return fail(400, { errors })
        }

        try {
            await locals.pb.collection('users').authWithPassword(email, password);
            if (!locals.pb?.authStore?.record) {
                locals.pb.authStore.clear();
                return {
                    notVerified: true
                };
            }
            // eslint-disable-next-line @typescript-eslint/no-explicit-any
        } catch (err: any) {
            console.log('Error: ', err);
            throw error(err.status, err.message);
        }

        throw redirect(303, '/account');

    },
}
