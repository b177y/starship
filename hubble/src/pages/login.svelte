<script>
	import { LoggedIn, AuthToken } from '../store.js';
    import { getNotificationsContext } from 'svelte-notifications';
    const { addNotification } = getNotificationsContext();
    import { navigate } from "svelte-routing";
    import API from '../api.js';
    let username = "";
    let password = "";
    export const loginRequest = async() => {
        try {
            const response = await API.post("/login",
                {
                    'username': username,
                    'password': password,
                })
            console.log(response)
            addNotification({
                text: `Signed In`,
                position: 'top-right',
                type: 'success',
                removeAfter: 2000,
            })
            LoggedIn.update(_ => true)
            AuthToken.update(_ => response.token)
            navigate('/')
        } catch(error){
            console.log(error);
            addNotification({
                text: `Could not sign in`,
                position: 'top-right',
                type: 'danger',
                removeAfter: 2000,
            })
        }
    }
</script>
<main class="p-4 flex justify-items-center w-full justify-center align-center">
    <div class="flex flex-col justify-center w-1/3 shadow p-8 mt-16">
        <h1 class="text-5xl">Login</h1>
        <div>
            <div class="flex flex-row mt-4">
                <label for="username" class="my-auto">Username</label>
                <input type="text" id="username" name="username"
                    class="w-1/2 h-6 pl-3 pr-6 text-base placeholder-gray-600 border rounded-lg appearance-none focus:shadow-outline my-auto ml-4 mr-2 bg-gray-100"
                    bind:value={username}
                >
            </div>
            <div class="flex flex-row mt-4 sm:bg-black md:bg-green-400 xl:bg-red-500">
                <label for="password" class="my-auto">Password</label>
                <input type="password" id="password" name="password"
                    class="w-1/2 h-6 pl-3 pr-6 text-base placeholder-gray-600 border rounded-lg appearance-none focus:shadow-outline my-auto ml-4 mr-2 bg-gray-100"
                    bind:value={password}
                >
            </div>
            <div class="flex justify-end">
                <button class="bg-green-500 text-white rounded-xl p-1 px-2"
                    on:click|preventDefault={async() => await loginRequest()}>Login</button>
            </div>
        </div>
    </div>
</main>
