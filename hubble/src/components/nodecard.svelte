<script>
    import NodeSettings from './nodesettings.svelte';
    export let groups;
    export let netname;
    export let node = {};
    export let statusColour = "bg-gray-300"
    let collapsed = true;
    $: if (node.status === "pending"){
        statusColour = "bg-gray-400";
    } else if (node.status === "active"){
        statusColour = "bg-green-400";
    } else {
        statusColour = "bg-red-400"
    }
    import API from '../api.js';
    import { createEventDispatcher } from 'svelte';
    const dispatch = createEventDispatcher();
    import { getNotificationsContext } from 'svelte-notifications';
    const { addNotification } = getNotificationsContext();
    export const approveNode = async() => {
        try {
            const response = await API.post(`/networks/${netname}/nodes/${node.name}/approve`);
            console.log(response)
            dispatch('updateNodes');
            addNotification({
                text: `Enabled node ${node.name}`,
                position: 'top-right',
                type: 'success',
                removeAfter: 2000,
            })
            return response;
        } catch (error) {
            console.error(error);
            addNotification({
                text: `Could not enable node ${node.name}`,
                position: 'top-right',
                type: 'danger',
                removeAfter: 2000,
            })
        }
    };
    export const disableNode = async() => {
        try {
            const response = await API.post(`/networks/${netname}/nodes/${node.name}/disable`);
            console.log(response)
            dispatch('updateNodes');
            addNotification({
                text: `Disabled node ${node.name}`,
                position: 'top-right',
                type: 'success',
                removeAfter: 2000,
            })
            return response;
        } catch (error) {
            console.error(error);
            addNotification({
                text: `Could not disable node ${node.name}`,
                position: 'top-right',
                type: 'danger',
                removeAfter: 2000,
            })
        }
    };
</script>
<main>
    <div class="shadow mt-2 p-4 hover:shadow-lg">
        <div class="flex flex-row justify-between">
        <div>
            <h3 class="text-gray-700">
                <b class="text-black">{node.name}</b>
                <span class="dot {statusColour}"></span>
                { node.latest_fetch == "NEVER" ? "Config not yet fetched by node." : `Latest fetch at ${node.latest_fetch}`}
            </h3>
            <h4 class="text-gray-400">{node.pubkey}</h4>
            <h4>Hostname: {node.hostname}</h4>
            <h4>Address: {node.address}</h4>
        </div>
        <div>
            {#if node.status === "pending"}
                <button on:click={async() => {await approveNode()}}
                    class="bg-green-600 rounded-xl p-1 text-white">
                    Approve
                </button>
            {:else if node.status === "active"}
                <button on:click={async() => {await disableNode()}}
                class="bg-red-500 rounded-xl p-1 text-white">Disable</button>
            {:else}
                <button on:click={async() => {await approveNode()}}
                class="bg-green-600 rounded-xl p-1 text-white">Enable</button>
            {/if}
        </div>
    </div>
        {#if !collapsed}
            <NodeSettings nodename={node.name} netname={netname} netgroups={groups}/>
        {/if}
        <div class="flex justify-end">
                    <i class="fas fa-times"></i>
            <button on:click={()=>{collapsed = !collapsed}} class="p-1 border-1 rounded-xl">
                { collapsed ? "EXPAND" : "COLLAPSE" }
            </button>
        </div>
    </div>
</main>

<style>
    .dot {
      height: 10px;
      width: 10px;
      border-radius: 50%;
      display: inline-block;
    }
</style>
