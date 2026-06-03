 export function formatTime(s: string) {
   if (!s) return ''
   const d = new Date(s)
   if (isNaN(d.getTime())) return s
   const pad = (n: number) => String(n).padStart(2, '0')
   return (
     `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ` +
     `${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`
   )
 }
