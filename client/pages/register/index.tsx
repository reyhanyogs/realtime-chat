import { useState } from 'react'
import { useRouter } from 'next/router'
import { UserInfo } from '../../modules/auth_provider'
import { API_URL } from '../../constants'

const Index = () => {
    const [username, setUsername] = useState('')
    const [email, setEmail] = useState('')
    const [password, setPassword] = useState('')

    const router = useRouter()

    const submitHandler =  (e: React.SyntheticEvent) => {
      e.preventDefault()
            const res = fetch(`${API_URL}/signup`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ username, email, password }),
              })
              .then((res => {
                return router.push('/login')
              }))
            .catch((err) => {
              console.log(err)
            })
    }

    const backHandler = (e: React.SyntheticEvent) => {
        router.push('/login')
    }

  return (
    <div className='flex items-center justify-center min-w-full min-h-screen'>
      <form className='flex flex-col md:w-1/5'>
        <div className='text-3xl font-bold text-center'>
          <span className='text-blue'>Register</span>
        </div>
        <input
          placeholder='username'
          className='p-3 mt-8 rounded-md border-2 border-grey focus:outline-none focus:border-blue'
          value={username}
          onChange={(e) => setUsername(e.target.value)}
        />
        <input
          placeholder='email'
          className='p-3 mt-4 rounded-md border-2 border-grey focus:outline-none focus:border-blue'
          value={email}
          onChange={(e) => setEmail(e.target.value)}
        />
        <input
          type='password'
          placeholder='password'
          className='p-3 mt-4 rounded-md border-2 border-grey focus:outline-none focus:border-blue'
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />
        <button
          className='p-3 mt-6 rounded-md bg-blue font-bold text-white'
          type='submit'
          onClick={submitHandler}
        >
          Register
        </button>
        <div className='mt-2'>
          <span className='text-blue cursor-pointer' onClick={backHandler}>
            Back
          </span>
        </div>
      </form>
    </div>
  )
}

export default Index
