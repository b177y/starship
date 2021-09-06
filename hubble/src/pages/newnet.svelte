<script>
    let name = "";
    let cidr = "";
    import API from '../api.js';
    import { navigate } from "svelte-routing";
    import { getNotificationsContext } from 'svelte-notifications';
    const { addNotification } = getNotificationsContext();
    import { createEventDispatcher } from 'svelte';
    const dispatch = createEventDispatcher();
    export const newNetwork = async() => {
        console.log("starting req")
        try {
            const response = await API.post("/networks/new",
                {
                    'name': name,
                    'cidr': cidr,
                })
            console.log(response)
            dispatch('updateNetworks');
            navigate(`/networks/${name}`)
            console.log("ADDING NOTIFICATION")
            addNotification({
                text: `Created network ${name}`,
                position: 'top-right',
                type: 'success',
                removeAfter: 2000,
            })
        } catch(error){
            console.log(error);
            console.log("ADDING NOTIFICATION")
            addNotification({
                text: `Could not create network ${name}`,
                position: 'top-right',
                type: 'danger',
                removeAfter: 2000,
            })
        }
    }
    function handleClick(){
        console.log("starting handle click")
        newNetwork()
    }
</script>

<main class="p-4 flex justify-items-center w-full justify-center align-center">
    <div class="flex flex-col justify-center w-8/12 shadow p-8 mt-8">
        <h1 class="text-5xl">New Network</h1>
        <div>
            <div class="flex flex-row mt-4">
                <label for="name" class="my-auto">Network Name</label>
                <input type="text" id="name" name="name"
                    class="w-1/4 h-6 pl-3 pr-6 text-base placeholder-gray-600 border rounded-lg appearance-none focus:shadow-outline my-auto ml-4 mr-2 bg-gray-100"
                    bind:value={name}>
            </div>
            <div class="flex flex-row mt-4">
                <label for="cidr" class="my-auto">CIDR</label>
                <input type="text" id="cidr" name="cidr"
                    class="w-1/4 h-6 pl-3 pr-6 text-base placeholder-gray-600 border rounded-lg appearance-none focus:shadow-outline my-auto ml-4 mr-2 bg-gray-100"
                    bind:value={cidr}>
            </div>
            <div class="flex justify-end">
                <button class="bg-green-500 text-white rounded-xl p-1 px-2"
                    on:click|preventDefault={handleClick}>create</button>
            </div>
        </div>
    </div>
</main>

