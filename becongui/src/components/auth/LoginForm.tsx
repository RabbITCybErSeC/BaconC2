import React, { useState, type FormEvent } from 'react';
import { ArrowRight } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import { login } from './../../services/authService.tsx';

interface LoginFormProps {
  onSubmitSuccess: () => void;
  onSubmitError: (errors: Error[]) => void;
}

const LoginForm: React.FC<LoginFormProps> = ({ onSubmitSuccess, onSubmitError }) => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const navigate = useNavigate();

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    if (!username || !password) {
      onSubmitError([new Error('Username and password are required')]);
      return;
    }

    setIsLoading(true);
    onSubmitError([]);

    try {
      await login({ username, password }); // Use the login function from AuthenticationService
      onSubmitSuccess();
      navigate('/agents'); // Redirect to agents after successful login
    } catch (error) {
      console.error('Login error:', error);

      let errorsToShow = [new Error('Login failed. Please check your credentials.')];

      if (error instanceof Error) {
        errorsToShow = [error];
      }

      onSubmitError(errorsToShow);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="px-5 py-7">
      <label htmlFor="username" className="font-semibold text-sm text-gray-600 dark:text-gray-200 pb-1 block">
        Username
      </label>
      <input
        id="username"
        name="username"
        type="text"
        required
        value={username}
        onChange={(e) => setUsername(e.target.value)}
        disabled={isLoading}
        className="border rounded-lg px-3 py-2 mt-1 mb-5 text-sm w-full dark:bg-slate-700 dark:text-gray-100 dark:border-slate-600 focus:outline-none focus:ring-2 focus:ring-blue-500"
      />
      <label htmlFor="password" className="font-semibold text-sm text-gray-600 dark:text-gray-200 pb-1 block">
        Password
      </label>
      <input
        id="password"
        name="password"
        type="password"
        required
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        disabled={isLoading}
        className="border rounded-lg px-3 py-2 mt-1 mb-5 text-sm w-full dark:bg-slate-700 dark:text-gray-100 dark:border-slate-600 focus:outline-none focus:ring-2 focus:ring-blue-500"
      />
      <button
        type="submit"
        disabled={isLoading}
        className="transition duration-200 bg-blue-500 hover:bg-blue-600 focus:bg-blue-700 focus:shadow-sm focus:ring-4 focus:ring-blue-500 focus:ring-opacity-50 text-white w-full py-2.5 rounded-lg text-sm shadow-sm hover:shadow-md font-semibold text-center inline-flex items-center justify-center disabled:opacity-50 disabled:cursor-not-allowed"
      >
        <span className="inline-block mr-2">{isLoading ? 'Logging in...' : 'Login'}</span>
        {!isLoading && <ArrowRight className="w-4 h-4 inline-block" />}
        {isLoading && (
          <svg
            className="animate-spin h-4 w-4 inline-block text-white"
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
          >
            <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
            <path
              className="opacity-75"
              fill="currentColor"
              d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
            ></path>
          </svg>
        )}
      </button>
    </form>
  );
};

export default LoginForm;
